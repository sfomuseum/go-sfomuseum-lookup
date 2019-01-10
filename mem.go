package lookup

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"io"
)

type MemLookup struct {
	Lookup
	table map[string]interface{}
}

func NewMemLookup() (Lookup, error) {
	table := make(map[string]interface{})
	return NewMemLookupWithTable(table)
}

func NewMemLookupWithCSVFromReader(fh io.Reader, lookup_key string, target_key string) (Lookup, error) {

	reader, err := csv.NewDictReader(fh)

	if err != nil {
		return nil, err
	}

	return NewMemLookupWithCSVReader(reader, lookup_key, target_key)
}

func NewMemLookupWithCSVReader(reader *csv.DictReader, lookup_key string, target_key string) (Lookup, error) {

	table := make(map[string]interface{})

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		k, ok := row[lookup_key]

		if !ok {
			msg := fmt.Sprintf("row is missing %s (lookup) key", lookup_key)
			return nil, errors.New(msg)
		}

		v, ok := row[target_key]

		if !ok {
			msg := fmt.Sprintf("row is missing %s (target) key", target_key)
			return nil, errors.New(msg)
		}

		table[k] = v
	}

	return NewMemLookupWithTable(table)
}

func NewMemLookupWithTable(table map[string]interface{}) (Lookup, error) {

	l := MemLookup{
		table: table,
	}

	return &l, nil
}

func (l *MemLookup) Find(path string) (interface{}, bool) {
	v, ok := l.table[path]
	return v, ok
}
