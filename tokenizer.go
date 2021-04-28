package main

import (
	"regexp"
	"strings"

	englishStemmer "github.com/kljensen/snowball/english"
)

type tokenizer struct {
	stopWords        map[string]struct{}
	punctuationRegex *regexp.Regexp
}

func newTokenizer() *tokenizer {
	stopWords := []string{
		"the", "be", "to", "of", "and", "a", "in", "that", "have",
		"I", "it", "for", "not", "on", "with", "he", "as", "you",
		"do", "at", "this", "but", "his", "by", "from", "wikipedia",
	}
	stopWordsSet := make(map[string]struct{})
	for _, word := range stopWords {
		stopWordsSet[word] = struct{}{}
	}
	re := regexp.MustCompile(`[^\w]`)
	return &tokenizer{
		stopWords:        stopWordsSet,
		punctuationRegex: re,
	}
}

func (t *tokenizer) tokenizeSingelWord(word string) string {
	// lowercase each token
	word = strings.ToLower(word)

	// remove any punctuation
	word = t.punctuationRegex.ReplaceAllString(word, "")
	if word == "" {
		return ""
	}

	// filter out if stopword
	if _, ok := t.stopWords[word]; ok {
		return ""
	}

	// apply stemming
	stemmed := englishStemmer.Stem(word, true)
	return stemmed
}

func (t *tokenizer) tokenize(text string) []string {
	var tokens []string
	// split text on white space and then tokenize
	for _, word := range strings.Fields(text) {
		tokenized := t.tokenizeSingelWord(word)
		if tokenized != "" {
			tokens = append(tokens, tokenized)
		}
	}
	return tokens
}
