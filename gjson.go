package lookup

import (
	"github.com/aaronland/go-storage"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"os"
)

type GJSONLookup struct {
	Lookup
	body []byte
}

func NewGJSONLookupFromStore(store storage.Store, path string) (Lookup, error) {

	fh, err := store.Get(path)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	return NewGJSONLookup(fh)
}

// deprecated - please just use NewGJSONLookupFromStore

func NewGJSONLookupFromFile(path string) (Lookup, error) {

	fh, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	return NewGJSONLookup(fh)
}

func NewGJSONLookup(fh io.Reader) (Lookup, error) {

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	l := GJSONLookup{
		body: body,
	}

	return &l, nil
}

func (l *GJSONLookup) Find(path string) (interface{}, bool) {

	rsp := gjson.GetBytes(l.body, path)

	if !rsp.Exists() {
		return nil, false
	}

	return rsp.Int(), true
}
