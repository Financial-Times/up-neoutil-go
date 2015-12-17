package neoutil

import (
	"github.com/jmcvetta/neoism"
)

type CypherRunner interface {
	WriteCypher(queries []*neoism.CypherQuery) error
}
