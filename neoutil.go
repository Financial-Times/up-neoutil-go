package neoutil

import (
	"github.com/Financial-Times/neoism"
)

// EnsureIndex creates an index on a label with a given property name if it does not already exist
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

// EnsureIndexes creates indexes for labels by the given property names, if they do not already exist
func EnsureIndexes(db *neoism.Database, labelToProperty map[string]string) error {
	for lab, prop := range labelToProperty {
		if err := EnsureIndex(db, lab, prop); err != nil {
			return err
		}
	}
	return nil
}
