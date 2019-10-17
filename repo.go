package lookup

import (
	"context"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"	
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"gocloud.dev/blob"
	"io"
	_ "log"
	"os"
	"strings"
	"sync"
)

type RepoLookup struct {
	Lookup
	lookup *sync.Map
}

func NewRepoLookupFromBucket(bucket *blob.Bucket, lookup_key string, target_key string) (Lookup, error) {

	lookup := new(sync.Map)

	iter := bucket.List(nil)
	ctx := context.Background()

	for {
		obj, err := iter.Next(ctx)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		path := obj.Key

		if !strings.HasSuffix(path, ".geojson") {
			continue
		}

		fh, err := bucket.NewReader(ctx, path, nil)

		if err != nil {
			return nil, err
		}

		defer fh.Close()
		
		f, err := feature.LoadGeoJSONFeatureFromReader(fh)

		if err != nil {
			return nil, err
		}

		recordFeature(f, lookup, lookup_key, target_key)

		// be explicit about this here since the defer above doesn't
		// get invoked as part of the for loop and it's too soon for
		// adding go routines and explcit throttling or queueing
		// (20191017/thisisaaronland)
		
		fh.Close()		
	}

	l := RepoLookup{
		lookup: lookup,
	}

	return &l, nil
}

func NewRepoLookupFromPath(root string, lookup_key string, target_key string) (Lookup, error) {

	lookup := new(sync.Map)

	cb := func(path string, info os.FileInfo) error {

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".geojson") {
			return nil
		}

		f, err := feature.LoadGeoJSONFeatureFromFile(path)

		if err != nil {
			// log.Println(err)
			return nil
		}

		recordFeature(f, lookup, lookup_key, target_key)
		
		return nil
	}

	c := crawl.NewCrawler(root)
	err := c.Crawl(cb)

	if err != nil {
		return nil, err
	}

	l := RepoLookup{
		lookup: lookup,
	}

	return &l, nil
}

func recordFeature(f geojson.Feature, lookup *sync.Map, lookup_key string, target_key string) error {

	lookup_rsp := gjson.GetBytes(f.Bytes(), lookup_key)

	if !lookup_rsp.Exists() {
		return errors.New("Missing lookup key")
	}

	target_rsp := gjson.GetBytes(f.Bytes(), target_key)

	if !target_rsp.Exists() {
		return errors.New("Missing target key")
	}

	k := lookup_rsp.String()
	v := target_rsp.String()

	// log.Println("SET", k, v)
	lookup.Store(k, v)

	return nil
}

func (l *RepoLookup) Find(path string) (interface{}, bool) {

	v, ok := l.lookup.Load(path)

	if !ok {
		return nil, false
	}

	return v.(string), true
}
