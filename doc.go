package main

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

type doc struct {
	ID       uint32
	Title    string
	URL      string
	Abstract string
}

func (d *doc) fullText() string {
	return strings.Join([]string{d.Title, d.Abstract}, " ")
}

func parseDocsFromXML(r io.Reader) ([]*doc, error) {
	var errNotStartElement = errors.New("not an xml start element")

	getElementData := func(decoder *xml.Decoder, t xml.Token) (key string, val string, err error) {
		elem, ok := t.(xml.StartElement)
		if !ok || elem.Name.Local == "doc" {
			return "", "", errNotStartElement
		}
		key = elem.Name.Local
		err = decoder.DecodeElement(&val, &elem)
		return
	}

	unflattenToDocs := func(kvs map[string][]string) []*doc {
		total := len(kvs["title"])
		var docs []*doc
		for i := 0; i < total; i++ {
			var doc doc
			doc.Title = kvs["title"][i]
			doc.URL = kvs["url"][i]
			doc.Abstract = kvs["abstract"][i]
			docs = append(docs, &doc)
		}
		return docs
	}

	decoder := xml.NewDecoder(r)
	kvs := make(map[string][]string)
	var err error
	for {
		var t xml.Token
		t, err = decoder.Token()
		if err != nil {
			break
		}

		var key, val string
		key, val, err = getElementData(decoder, t)
		if err == errNotStartElement {
			continue
		}
		if err != nil {
			break
		}
		kvs[key] = append(kvs[key], val)
	}
	if err != nil && err != io.EOF {
		return nil, err
	}
	docs := unflattenToDocs(kvs)
	return docs, nil
}
