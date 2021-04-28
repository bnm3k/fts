package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/RoaringBitmap/roaring"
)

type counter map[string]int

type searchStore struct {
	tokenizer       *tokenizer
	documents       map[uint32]*doc
	tokenIndex      map[string]*roaring.Bitmap
	termFrequencies map[uint32]counter
	currID          uint32
}

func newSearchStore() *searchStore {
	return &searchStore{
		tokenizer:       newTokenizer(),
		documents:       make(map[uint32]*doc),
		tokenIndex:      make(map[string]*roaring.Bitmap),
		termFrequencies: make(map[uint32]counter),
		currID:          1,
	}
}

func (s *searchStore) insertDocument(doc *doc) uint32 {
	// assign unique ID
	id := s.currID
	doc.ID = id
	s.currID++

	// store document
	s.documents[id] = doc

	// index tokens
	tokens := s.tokenizer.tokenize(doc.fullText())
	for _, token := range tokens {
		set, ok := s.tokenIndex[token]
		if !ok {
			set = roaring.NewBitmap()
			s.tokenIndex[token] = set
		}
		set.Add(id)
	}

	// term frequencies for ranking
	counter := make(map[string]int)
	for _, token := range tokens {
		if _, ok := counter[token]; !ok {
			counter[token] = 1
		} else {
			counter[token]++
		}
	}
	s.termFrequencies[id] = counter

	return doc.ID
}

type searchResult struct {
	doc   *doc
	score float64
}

func (sr searchResult) PrettyString() string {
	abstract := sr.doc.Abstract
	if len(abstract) > 100 {
		abstract = abstract[:100] + "..."
	}
	return fmt.Sprintf("Title: %s\nScore: %f\nAbstract: %s\n",
		sr.doc.Title, sr.score, abstract)
}

func (s *searchStore) search(query string) []searchResult {
	queryTokens := s.tokenizer.tokenize(query)
	bitmaps := make([]*roaring.Bitmap, 0, len(queryTokens))
	for _, token := range queryTokens {
		idBitmap := s.searchSingleWord(token, true)
		if idBitmap != nil {
			bitmaps = append(bitmaps, idBitmap)
		}
	}
	intersection := roaring.FastAnd(bitmaps...)

	// retrieve and score documents based on result IDs
	var results []searchResult
	iter := intersection.ManyIterator()
	chunkOfIDs := make([]uint32, 128)
	for {
		n := iter.NextMany(chunkOfIDs)
		if n == 0 {
			break
		}
		for _, id := range chunkOfIDs[:n] {
			doc := s.documents[id]

			// calculate score for relevancy
			termFrequencies := s.termFrequencies[id]
			score := float64(0)
			for _, token := range queryTokens {
				tf, ok := termFrequencies[token]
				idf := math.Log10(float64(len(s.documents)) / float64(s.tokenIndex[token].GetCardinality()))
				if ok {
					score += (float64(tf) * idf)
				}
			}
			results = append(results, searchResult{
				doc:   doc,
				score: score,
			})
		}
	}

	// sort results based on score
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	return results
}

func (s *searchStore) searchSingleWord(word string, tokenized bool) *roaring.Bitmap {
	if !tokenized {
		word = s.tokenizer.tokenizeSingelWord(word)
	}
	return s.tokenIndex[word]
}
