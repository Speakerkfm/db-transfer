package main

import (
	"crypto/rand"
	"crypto/rsa"
	flagsCfg "db-transfer/consumer/pkg/flags"
	"db-transfer/consumer/pkg/log"
	"db-transfer/consumer/pkg/models"
	"db-transfer/consumer/pkg/services"
	"db-transfer/consumer/pkg/store"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"
	"github.com/jessevdk/go-flags"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const programName = "data-transfer-consumer"

var (
	envFlags   *flagsCfg.Config
	parser     *flags.Parser
	messageHandler   *services.MessageHandler
	db         *gorm.DB
	mqConsumer *services.MqConsumer
	tcpServer  *services.TcpServer
)

func init() {
	envFlags = &flagsCfg.Config{}

	parser = flags.NewParser(envFlags, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	log.Config(programName, os.Stderr)

	//mysql
	mysqlConf := mysql.NewConfig()
	mysqlConf.Net = "tcp"
	mysqlConf.Addr = fmt.Sprintf("%s:%v", envFlags.DatabaseHost, envFlags.DatabasePort)
	mysqlConf.User = envFlags.DatabaseUser
	mysqlConf.Passwd = envFlags.DatabasePassword
	mysqlConf.DBName = envFlags.DatabaseName
	mysqlConf.ParseTime = true
	mysqlConf.Loc = time.Local

	var err error
	db, err = gorm.Open("mysql", mysqlConf.FormatDSN())
	if err != nil {
		log.Fatal().Err(err).Msg("Mysql connection failed")
	}

	db.SingularTable(true)

	//log sql queries
	db.LogMode(false)

	//redis
	redisOpt := redis.Options{Addr: envFlags.RedisHost}
	redisClient := redis.NewClient(&redisOpt)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	st := store.NewStore(db, redisClient)

	transfer := services.NewTransfer(st)
	sessionsService := services.NewSessionService(st)
	messageHandler = services.NewMessageHandler(transfer, sessionsService, st, privateKey)
}

func main() {
	defer func() {
		if err := db.Close(); err != nil {
			log.Warn().Err(err).Msg("db close error")
		}
		log.Info().Msg(fmt.Sprintf("mysql connect close"))
	}()

	mqConsumer = services.NewMqConsumer(envFlags, consumeFunc)

	tcpServer = services.NewTcpServer(envFlags, handlerFunc)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-sc
	log.Info().Msg(fmt.Sprintf("shutting down by signal %s", sig))

	if err := mqConsumer.Shutdown(); err != nil {
		log.Fatal().Err(err).Msg("error during shutdown")
	}

	if err := tcpServer.Shutdown(); err != nil {
		log.Fatal().Err(err).Msg("error during shutdown")
	}
}

func consumeFunc(message amqp.Delivery) {
	var msg *models.Message

	if err := json.Unmarshal(message.Body, &msg); err != nil {
		log.Warn().Interface("body", message.Body).Msg("Bad message format")
		if err = message.Nack(false, false); err != nil {
			log.Warn().Err(err).Msg("Message nack error")
		}
		return
	}

	defer recoverFunc(message)

	if _, err := messageHandler.Handle(msg); err != nil {
		log.Error().Err(err).Msg("Failed to handle message")

		if err := message.Nack(false, false); err != nil {
			log.Warn().Err(err).Msg("Message nack error")
		}
		return
	}

	if err := message.Ack(false); err != nil {
		log.Warn().Err(err).Msg("Message ack error")

		return
	}
}

func recoverFunc(message amqp.Delivery) {
	if r := recover(); r != nil {
		if err, _ := r.(error); err != nil {
			log.Warn().Err(fmt.Errorf("pkg: %v", r)).Msg("Message handle error")

			if err := message.Nack(false, false); err != nil {
				log.Warn().Err(err).Msg("Message ack error")
			}
		}
	}
}

func handlerFunc(messageBody []byte) []byte{
	var msg *models.Message

	if err := json.Unmarshal(messageBody, &msg); err != nil {
		log.Warn().Interface("body", messageBody).Msg("Bad message format")

		return nil
	}

	respMsg, err := messageHandler.Handle(msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to handle message")

		return nil
	}

	messageBody, err = json.Marshal(respMsg)
	if err != nil {
		log.Warn().Msg("Bad response message format")

		return nil
	}

	return messageBody
}
