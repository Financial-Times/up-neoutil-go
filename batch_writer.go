package neoutil

import (
	"github.com/Financial-Times/neo-cypher-runner-go"
	"github.com/jmcvetta/neoism"
	"time"
)

func NewBatchWriter(db *neoism.Database, batchSize int) CypherRunner {
	return neocypherrunner.NewBatchCypherRunner(db, batchSize, time.Millisecond*20)
}

type CypherRunner interface {
	CypherBatch(queries []*neoism.CypherQuery) error
}
