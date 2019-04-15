package lookup

import (
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	_ "log"
	"os"
	"strings"
	"sync"
)

type RepoLookup struct {
	Lookup
	lookup *sync.Map
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

		lookup_rsp := gjson.GetBytes(f.Bytes(), lookup_key)

		if !lookup_rsp.Exists() {
			return nil
		}

		target_rsp := gjson.GetBytes(f.Bytes(), target_key)

		if !target_rsp.Exists() {
			return nil
		}

		k := lookup_rsp.String()
		v := target_rsp.String()

		// log.Println("SET", k, v)
		lookup.Store(k, v)
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

func (l *RepoLookup) Find(path string) (interface{}, bool) {

	v, ok := l.lookup.Load(path)

	if !ok {
		return nil, false
	}

	return v.(string), true
}
