package bitcoin

import (
	"encoding/json"

	// pq will bind to database/sql
	_ "github.com/lib/pq"
)

// WriteBlockToDB comment
func (b *Bitcoind) WriteBlockToDB(block Block) (err error) {

	blockJSON, _ := json.Marshal(block)

	txJSON, _ := json.Marshal(block.Tx)

	tx, _ := b.GetRawTransaction(block.Tx[0], false)
	cbJSON, _ := json.Marshal(tx)

	insertStmt, err := b.db.Prepare(`
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
		b.coin,
		block.Hash,
		block.Tx[0],
		block.Size,
		len(block.Tx),
		block.Height,
		block.Difficulty,
		0,
		string(cbJSON),
		string(txJSON),
		string(blockJSON),
	)

	return
}
