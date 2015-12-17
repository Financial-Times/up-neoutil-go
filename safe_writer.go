package neoutil

import (
	"github.com/jmcvetta/neoism"
	"log"
	"time"
)

// SafeWriter does not ensure writes complete before returning
type safeWriter struct {
	db         *neoism.Database
	writeQueue chan writeEntry
	batchSize  int
}

type writeEntry struct {
	queries []*neoism.CypherQuery
	err     chan error
}

// NewSafeWriter provides a new batch writer which will batch writes internally
// without risking data loss.
func NewSafeWriter(db *neoism.Database, batchSize int) *safeWriter {
	sw := &safeWriter{db, make(chan writeEntry, batchSize), batchSize}
	go sw.writeLoop()
	return sw
}

func (sw *safeWriter) WriteCypher(queries []*neoism.CypherQuery) error {
	we := writeEntry{queries, make(chan error)}
	sw.writeQueue <- we
	return <-we.err
}

func (sw *safeWriter) writeLoop() {

	var qs []writeEntry

	timer := time.NewTimer(1 * time.Second)

	defer log.Println("write loop exited")
	for {
		select {
		case writeEntry, ok := <-sw.writeQueue:
			if !ok {
				return
			}
			qs = append(qs, writeEntry)

			if len(qs) < sw.batchSize {
				timer.Reset(20 * time.Millisecond)
				continue
			}
		case <-timer.C:
		}

		if len(qs) > 0 {
			timer.Stop()
			log.Printf("writing batch of %d\n", len(qs))
			var batched []*neoism.CypherQuery
			var errChs []chan error
			for _, we := range qs {
				batched = append(batched, we.queries...)
				errChs = append(errChs, we.err)
			}

			err := sw.db.CypherBatch(batched)

			/*  not batched
			var err error
			for _, q := range batched {
				if err == nil {
					err = sw.db.Cypher(q)
				}
			}
			*/

			// transactional
			//tx, err := sw.db.Begin(batched)
			//if err == nil {
			//		err = tx.Commit()
			//	}

			for _, errCh := range errChs {
				errCh <- err
			}

			log.Printf("wrote batch of %d\n", len(qs))
			qs = qs[0:0]
		}
	}
}
