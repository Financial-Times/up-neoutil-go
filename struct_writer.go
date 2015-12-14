package neoutil

import (
	"github.com/jmcvetta/neoism"
)

// deprecated
type CypherWriter interface {
	WriteCypher(queries []*neoism.CypherQuery) error
	Close() error
}

type CypherRunner interface {
	WriteCypher(queries []*neoism.CypherQuery) error
}
