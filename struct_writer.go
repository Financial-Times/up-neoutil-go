package neoutil

import (
	"github.com/jmcvetta/neoism"
)

type CypherWriter interface {
	WriteCypher(queries []*neoism.CypherQuery) error
	Close() error
}
