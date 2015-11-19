package neoutil

import (
	"fmt"
	"github.com/jmcvetta/neoism"
	"log"
	"time"
)

// BatchWriter provides a way to batch up writes to neo.
type BatchWriter struct {
	db         *neoism.Database
	WriteQueue chan<- []*neoism.CypherQuery
	Closed     <-chan struct{}
}

// NewBatchWriter provides a new batch writer, which will flush writes either
// when there are at least 1024 queries, or when 1 second has passed without
// any new queries being queued, whichever happens first.
func NewBatchWriter(db *neoism.Database) *BatchWriter {
	wq := make(chan []*neoism.CypherQuery)

	closed := make(chan struct{})

	bw := &BatchWriter{db: db, WriteQueue: wq, Closed: closed}
	go bw.writeLoop(wq, closed)
	return bw
}

func (bw *BatchWriter) writeLoop(writeQueue <-chan []*neoism.CypherQuery, closed chan struct{}) {

	var qs []*neoism.CypherQuery

	timer := time.NewTimer(1 * time.Second)

	defer log.Println("write loop exited")
	defer close(closed)
	for {
		select {
		case o, ok := <-writeQueue:
			if !ok {
				return
			}
			for _, q := range o {
				qs = append(qs, q)
			}
			if len(qs) < 1024 {
				timer.Reset(1 * time.Second)
				continue
			}
		case <-timer.C:
		}
		if len(qs) > 0 {
			timer.Stop()
			fmt.Printf("writing batch of %d\n", len(qs))
			err := bw.db.CypherBatch(qs)
			if err != nil {
				panic(err)
			}
			fmt.Printf("wrote batch of %d\n", len(qs))
			qs = qs[0:0]
		}
	}
}
