package main

import (
	"flag"
	"github.com/sfomuseum/go-sfomuseum-lookup"
	"log"
)

func main() {

	lookup_key := flag.String("lookup", "", "...")
	target_key := flag.String("target", "", "...")

	flag.Parse()

	for _, path := range flag.Args() {

		l, err := lookup.NewRepoLookupFromPath(path, *lookup_key, *target_key)

		if err != nil {
			log.Fatal(err)
		}

		log.Println(l)
	}

}
