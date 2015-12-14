package neoutil

import (
	"encoding/json"
)

type NeoEngine interface {
	CreateOrUpdate(cr CypherRunner, obj interface{}) error
	Delete(cr CypherRunner, identity string) (deleted bool, err error)
	SuggestedIndexes() map[string]string
	DecodeJSON(*json.Decoder) (obj interface{}, identity string, err error)
	Read(cr CypherRunner, identity string) (obj interface{}, found bool, err error)
}
