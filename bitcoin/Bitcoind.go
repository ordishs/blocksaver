package bitcoin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ordishs/gocore"
	zmq "github.com/pebbe/zmq4"
)

// A Bitcoind represents a Bitcoind client
type Bitcoind struct {
	coin         string
	client       *rpcClient
	db           *sql.DB
	zmqConnected bool
}

// New return a new bitcoind
func New(coin string) (*Bitcoind, error) {
	var (
		username, _ = gocore.Config().Get(coin + "_user")
		password, _ = gocore.Config().Get(coin + "_password")
		host, _     = gocore.Config().Get(coin + "_host")
		port, _     = gocore.Config().GetInt(coin + "_port")
		useSSL      = false
		zmqPort, _  = gocore.Config().GetInt(coin + "_zmqPort")

		dbHost, _     = gocore.Config().Get("db_host")
		dbPort, _     = gocore.Config().GetInt("db_port")
		dbName, _     = gocore.Config().Get("db_name")
		dbUser, _     = gocore.Config().Get("db_user")
		dbPassword, _ = gocore.Config().Get("db_password")
	)

	rpcClient, err := newClient(host, port, username, password, useSSL)
	if err != nil {
		return nil, err
	}

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%d", dbUser, dbPassword, dbName, dbHost, dbPort)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}

	bitcoind := &Bitcoind{
		coin:   coin,
		client: rpcClient,
		db:     db,
	}

	if zmqPort != 0 {
		zmqAddress := fmt.Sprintf("tcp://%s:%d", host, zmqPort)
		bitcoind.connectToZMQ(zmqAddress)
	}

	return bitcoind, nil
}

// GetBlock returns information about the block with the given hash.
// berr is a bitcoin specific error
func (b *Bitcoind) GetBlock(blockHash string) (block Block, berr Error, err error) {
	verbose := true
	r, err := b.client.call("getblock", []interface{}{blockHash, verbose})

	if err != nil {
		if r.Err != nil {
			rr := r.Err.(map[string]interface{})
			berr = Error{
				Code:    rr["code"].(float64),
				Message: rr["message"].(string),
			}
		}
		return
	}

	if !verbose {
		fmt.Print(string(r.Result))
	} else {
		err = json.Unmarshal(r.Result, &block)
		if err != nil {
			return
		}
	}

	return
}

// GetRawTransaction returns raw transaction representation for given transaction id.
func (b *Bitcoind) GetRawTransaction(txID string, verbose bool) (rawTx RawTransaction, err error) {
	intVerbose := 0
	if verbose {
		intVerbose = 1
	}
	r, err := b.client.call("getrawtransaction", []interface{}{txID, intVerbose})
	if err != nil {
		return
	}
	if !verbose {
		err = json.Unmarshal(r.Result, &rawTx)
	} else {
		var t RawTransaction
		err = json.Unmarshal(r.Result, &t)
		rawTx = t
	}
	return
}

func (b *Bitcoind) connectToZMQ(zmqAddress string) {
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
				continue
			}

			if b.zmqConnected == false {
				b.zmqConnected = true
				log.Printf("ZMQ: Subscription to %s established\n", zmqAddress)
				subscriber.SetUnsubscribe("rawblock")
				subscriber.SetUnsubscribe("hashtx")
				subscriber.SetUnsubscribe("rawtx")
			}

			if msg == "hashblock" {
				fmt.Printf("We've got a block message, %s\n", msg)
				bytes, err := subscriber.RecvBytes(0)
				if err != nil {
					fmt.Printf("ERROR: %+v", err)
				} else {
					blockHash := fmt.Sprintf("%x", bytes)
					block, berr, err := b.GetBlock(blockHash)
					if err != nil || berr.Code != 0 {
						fmt.Printf("ERROR: %+v, %+v", err, berr)
					}
					fmt.Printf("We've got block message, %+v", block)
					// err = b.WriteBlockToDB(block)
				}
			} else {
				fmt.Printf("We've got a message, %+v\n", msg)
			}
			fmt.Printf("Message: %s\n%x\n\n", msg, msg)
		}
	}()

}
