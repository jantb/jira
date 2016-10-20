package main

import (
	"testing"
)
var data1 = `This is a sample a`
var data12 = `This is another another example example example`
var data13 = `This is another another example example where we are wrong`

func Test_cosine_similarity(t *testing.T) {
	tfidf := tfidf([]string{data1, data12,data13})

	if similarity(tfidf[data13], tfidf[data12]) != 0.4537233687954632 {
		t.Error("Expected similarity(tfidf[data], tfidf[data2]) to be 0.4537233687954632 was", similarity(tfidf[data12], tfidf[data13]))
	}

}