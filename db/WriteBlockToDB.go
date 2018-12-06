package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	// pq will bind to database/sql
	_ "github.com/lib/pq"

	"../bitcoin"
)

// WriteBlockToDB comment
func WriteBlockToDB(block bitcoin.Block) (err error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%d", dbUser, dbPassword, dbName, dbHost, dbPort)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return
	}

	defer db.Close()

	blockJSON, _ := json.Marshal(block)

	txJSON, _ := json.Marshal(block.Tx)

	tx, _ := bitcoin.GetRawTransaction(block.Tx[0])
	cbJSON, _ := json.Marshal(tx)

	insertStmt, err := db.Prepare(`
		INSERT INTO blocks (
		 Coin
		,Hash
		,cbHash
		,Size
		,txCount
		,Height
		,Difficulty
		,Reward
		,coinbaseJSON
		,txJSON
		,blockJSON
		) VALUES (
		 $1
		,$2
		,$3
		,$4
		,$5
		,$6
		,$7
		,$8
		,$9
		,$10
		,$11
		)`)

	_, err = insertStmt.Exec(
		bitcoin.Coin,
		block.Hash,
		block.Tx[0],
		block.Size,
		len(block.Tx),
		block.Height,
		block.Difficulty,
		0,
		string(cbJSON),
		string(txJSON),
		string(blockJSON))

	return
}
