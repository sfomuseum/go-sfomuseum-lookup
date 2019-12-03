package lookup

import (
	"context"
	"gocloud.dev/blob"
	"github.com/tidwall/gjson"
	"path/filepath"
	"io/ioutil"
)

type GJSONLookup struct {
	Lookup
	body []byte
}

func NewGJSONLookup() Lookup {
	lu := &GJSONLookup{}
	return lu
}

func (lu *GJSONLookup) Open(ctx context.Context, uri string) error {

	root := filepath.Dir(uri)
	fname := filepath.Base(uri)	
	
	bucket, err := blob.OpenBucket(ctx, root)

	if err != nil {
		return err
	}

	fh, err := bucket.NewReader(ctx, fname, nil)

	if err != nil {
		return err
	}
	
	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return err
	}

	lu.body = body
	return nil
}

func (l *GJSONLookup) Find(path string) (interface{}, bool) {

	rsp := gjson.GetBytes(l.body, path)

	if !rsp.Exists() {
		return nil, false
	}

	return rsp.Int(), true
}
