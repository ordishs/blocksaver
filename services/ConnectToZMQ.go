package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ordishs/gocore"
	zmq "github.com/pebbe/zmq4"

	"../bitcoin"
	"../db"
)

// ConnectToZMQ comment
func ConnectToZMQ() {
	host, _ := gocore.Config().Get(bitcoin.Coin + "_host")
	zmqPort, _ := gocore.Config().GetInt(bitcoin.Coin + "_zmqPort")
	zmqAddress := fmt.Sprintf("tcp://%s:%d", host, zmqPort)
	connected := false

	go func() {
		//  First, connect our subscriber socket
		subscriber, err := zmq.NewSocket(zmq.SUB)
		if err != nil {
			log.Fatal(err)
		}
		defer subscriber.Close()
		subscriber.Connect(zmqAddress)
		subscriber.SetSubscribe("hashblock")
		subscriber.SetSubscribe("hashtx")
		subscriber.SetSubscribe("rawblock")
		subscriber.SetSubscribe("rawtx")

		log.Printf("ZMQ: Subscribing to %s", zmqAddress)

		//  0MQ is so fast, we need to wait a while...
		time.Sleep(time.Second)
		for {
			msg, err := subscriber.Recv(0)
			if err != nil {
				log.Printf("ERROR: %+v", err)
			} else {
				if connected == false {
					connected = true
					log.Printf("ZMQ: Subscription to %s established\n", zmqAddress)
					subscriber.SetUnsubscribe("hashblock")
					subscriber.SetUnsubscribe("hashtx")
					subscriber.SetUnsubscribe("rawtx")
				}
				if msg == "rawblock" {
					fmt.Printf("We've got a raw block message, %+v", msg)
					b, err := subscriber.RecvBytes(0)
					if err != nil {
						fmt.Printf("ERROR: %+v", err)
					} else {
						var block bitcoin.Block
						json.Unmarshal(b, &block)
						fmt.Printf("We've got block message, %+v", block)
						err = db.WriteBlockToDB(block)
					}
				} else {
					fmt.Printf("We've got a message, %+v\n", msg)
				}
			}
		}
	}()

}
