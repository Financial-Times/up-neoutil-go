package neoutil

import (
	"github.com/Financial-Times/neo-cypher-runner-go"
	"github.com/jmcvetta/neoism"
)

func NewBatchWriter(db *neoism.Database, batchSize int) CypherRunner {
	//TODO: remove this pointless wrapper and update callers at some point.
	return neocypherrunner.NewBatchCypherRunner(db, batchSize)
}

type CypherRunner interface {
	CypherBatch(queries []*neoism.CypherQuery) error
}
