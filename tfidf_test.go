package main

import (
	"testing"
)

var data = `This is a sample a`
var data2 = `This is another another example example example`

func TestTfidf(t *testing.T) {
	tfidf := tfidf([]string{data, data2})
	if tfidf[data]["example"] != 0 {
		t.Error("Expected tfidf[data2][\"example\"] to be 0 was", tfidf[data]["example"])
	}
	if tfidf[data2]["example"] != 0.12901285528456335 {
		t.Error("Expected tfidf[data2][\"example\"] to be 0.12901285528456335 was", tfidf[data2]["example"])
	}

}
func TestTf(t *testing.T) {
	if tf(data)["this"] != 0.2 {
		t.Error("Expected tf(data)[\"this\"] to be 0.2 was", tf(data)["this"])
	}
}

func TestTf_example(t *testing.T) {
	if tf(data)["example"] != 0 {
		t.Error("Expected tf(data)[\"example\"] to be 00 was", tf(data)["example"])
	}
	if tf(data2)["example"] != 0.42857142857142855 {
		t.Error("Expected tf(data2)[\"example\"] to be 0.42857142857142855 was", tf(data2)["example"])
	}
}

func TestIdf(t *testing.T) {
	if idf([]string{data, data2})["this"] != 0 {
		t.Error("Expected idf(data)[\"this\"] to be 0 was", idf([]string{data, data2})["this"])
	}
}

func TestIdf_example(t *testing.T) {
	if idf([]string{data, data2})["example"] != 0.3010299956639812 {
		t.Error("idf([]string{data, data2})[\"example\"] to be 0 was", idf([]string{data, data2})["example"])
	}
}