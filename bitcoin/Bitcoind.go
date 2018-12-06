package bitcoin

import (
	"encoding/json"
	"fmt"

	"github.com/ordishs/gocore"
)

// A Bitcoind represents a Bitcoind client
type Bitcoind struct {
	client *rpcClient
}

// New return a new bitcoind
func New(coin string) (*Bitcoind, error) {
	var (
		username, _ = gocore.Config().Get(coin + "_user")
		password, _ = gocore.Config().Get(coin + "_password")
		host, _     = gocore.Config().Get(coin + "_host")
		port, _     = gocore.Config().GetInt(coin + "_port")
		useSSL      = false
	)

	rpcClient, err := newClient(host, port, username, password, useSSL)
	if err != nil {
		return nil, err
	}

	return &Bitcoind{rpcClient}, nil
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
