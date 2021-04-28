package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

func main() {
	buf, err := ioutil.ReadFile("./data/testdata.xml")
	if err != nil {
		panic(err)
	}

	r := bytes.NewBuffer(buf)
	docs, err := parseDocsFromXML(r)
	if err != nil {
		panic(errors.Wrap(err, "parse docs from XML"))
	}

	searchStore := newSearchStore()
	for _, d := range docs {
		searchStore.insertDocument(d)
	}

	searchResults := searchStore.search("foo")
	for _, res := range searchResults {
		fmt.Println(res.PrettyString())
	}
}
