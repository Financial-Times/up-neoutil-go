package neoutil

import (
	"github.com/jmcvetta/neoism"
)

func EnsureIndex(db *neoism.Database, label string, prop string) error {
	indexes, err := db.Indexes(label)
	if err != nil {
		return err
	}
	for _, ind := range indexes {
		if len(ind.PropertyKeys) == 1 && ind.PropertyKeys[0] == prop {
			return nil
		}
	}
	if _, err := db.CreateIndex(label, prop); err != nil {
		return err
	}
	return nil
}

func EnsureIndexes(db *neoism.Database, labelToProperty map[string]string) error {
	for lab, prop := range labelToProperty {
		if err := EnsureIndex(db, lab, prop); err != nil {
			return err
		}
	}
	return nil
}
