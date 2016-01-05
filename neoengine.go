package neoutil

import (
	"encoding/json"
)

type NeoEngine interface {
	CreateOrUpdate(obj interface{}) error
	Delete(identity string) (deleted bool, err error)
	SuggestedIndexes() map[string]string
	DecodeJSON(*json.Decoder) (obj interface{}, identity string, err error)
	Read(dentity string) (obj interface{}, found bool, err error)
}
