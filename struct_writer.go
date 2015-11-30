package neoutil

import (
	"github.com/jmcvetta/neoism"
)

type StructWriter interface {
	WriteCipher(queries []*neoism.CypherQuery) error
	Close() error
}
