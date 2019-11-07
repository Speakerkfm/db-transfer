package main

import (
	"database/sql"
	flagsCfg "db-transfer/producer/flags"
	"db-transfer/producer/pkg/service"
	"db-transfer/producer/pkg/store"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jessevdk/go-flags"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"os"
	"time"
)

var (
	envFlags *flagsCfg.Config
	parser   *flags.Parser
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
}

func main() {
	//sqlite
	db, err := sql.Open("sqlite3", "./test.db")
	checkErr(err)
	defer db.Close()

	//rabbit
	amqpClient, err := service.NewMQ(envFlags)
	if err != nil {
		panic(err)
	}
	defer amqpClient.Connection.Close()

	//redis
	redisOpt := redis.Options{Addr: envFlags.RedisHost}
	redisClient := redis.NewClient(&redisOpt)

	st := store.NewStore(db, redisClient)

	myConn := &service.Conn{
		IdleTimeout: 5 * time.Second,
		MaxReadBuffer: 262144,
	}

	clientID, _ := uuid.NewV4()

	client := service.NewTcpClient(envFlags.ServerHost, myConn)
	sessionService := service.NewSessionService(client, st)

	session, err := sessionService.CreateSession(clientID.String())
	checkErr(err)

	rawData, err := st.GetData()
	checkErr(err)

	data, err := json.Marshal(rawData)
	checkErr(err)

	encryptedData, err := sessionService.GetEncryptedDataAES(session.ClientID, session.Key, service.ImportData, data)
	checkErr(err)

	switch envFlags.SendType {
	case "tcp":
		//socket
		fmt.Println("Sending data with tcp connection")
		err = client.SendTcpRequest(encryptedData)
		checkErr(err)

		break
	case "mq":
		//amqp
		fmt.Println("Sending data with RabbitMQ")
		err = amqpClient.Publish(encryptedData)
		checkErr(err)

		break
	}

	fmt.Println("Done!")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
