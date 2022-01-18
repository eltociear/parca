package columnstore

import "fmt"

type Table struct {
	Name     string
	Schema   *Schema
	Granules []*Granule
}

func NewTable(name string, schema *Schema) *Table {
	return &Table{
		Name:     name,
		Schema:   schema,
		Granules: []*Granule{},
	}
}

// Insert inserts rows into the table. Rows are expected to already be sorted
// by the table schema's sort columns.
func (t *Table) Insert(rows ...Row) error {
	// Ensure that map column keys are a superset of the map keys of `rows`.

	// Split rows into batches grouped respectively by the granule they belong
	// to and insert them into that granule.

	g := NewGranule(t.Schema)
	_, err := g.Insert(rows...)
	if err != nil {
		return err
	}
	// TODO: If granule is full, compact and split.
	// TODO: Insert any new granule into the sparse index.
	t.Granules = append(t.Granules, g)

	return nil
}

// String prints out the table contents in human readable rows
func (t *Table) String() string {
	var s string
	for _, granule := range t.Granules {
		for _, part := range granule.Parts {
			for _, col := range part.Columns {

				switch t := col.(type) {
				case *PlainColumn:
					s += fmt.Sprint(t.def.Name)
					s += " := "
					s += fmt.Sprint(t.enc)
					s += "\n"
				case *MapColumn:
					s += fmt.Sprint(t.def.Name)
					s += "\n"
					for key, c := range t.columns {
						s += fmt.Sprint(key)
						s += ": "
						s += fmt.Sprint(c)
					}
				}
			}
		}
	}

	return s
}
