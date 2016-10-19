package main

import (
	"strings"
	"math"
	"fmt"
)

func toLower(words string) []string {
	return strings.Fields(strings.ToLower(words))
}

func removeNorwegianStopwords(words []string) []string {
	stopwords := strings.Fields(`å alle andre arbeid av begge bort bra bruke da denne der deres det din disse du eller en ene
	eneste enhver enn er et få folk for fordi forsøke fra før først gå gjorde gjøre god ha hadde han hans hennes
	her hva hvem hver hvilken hvis hvor hvordan hvorfor i ikke inn innen kan kunne lage lang lik like må makt mange
	måte med meg meget men mens mer mest min mye nå når navn nei ny og også om opp oss over på part punkt rett
	riktig så samme sant si siden sist skulle slik slutt som start stille tid til tilbake tilstand under ut uten
	var vår ved verdi vi vil ville vite være vært`)
	var ret []string
	for _, word := range words {
		found := false
		for _, stopword := range stopwords {
			if stopword == word {
				found = true
			}
		}
		if !found {
			ret = append(ret, word)
		}
	}

	return ret
}

func wf(words []string) (map[string]int) {
	m := make(map[string]int)
	for _, word := range words {
		if _, exists := m[word]; exists {
			m[word]++
		} else {
			m[word] = 1
		}
	}
	return m
}
func normalize(m map[string]int) map[string]float64 {
	var count int
	for _, value := range m {
		count += value
	}
	r := make(map[string]float64)
	for key, value := range m {
		r[key] = float64(value) / float64(count)
	}
	return r
}

func set(words []string) []string {
	var set []string
	for _, word := range words {
		contains := false
		for _, w := range set {
			if strings.EqualFold(w, word) {
				contains = true
				break
			}
		}
		if !contains {
			set = append(set, word)
		}
	}
	return set
}

func tf(string string) (map[string]float64) {
	return normalize(wf(removeNorwegianStopwords(toLower(string))))
}
func idf(strings []string) (map[string]float64) {
	d := float64(len(strings))
	var words []string
	for _, s := range strings {
		words = append(words, set(removeNorwegianStopwords(toLower(s)))...)
	}
	idfMap := make(map[string]float64)
	for word, value := range wf(words) {
		idfMap[word] = math.Log10(d / float64(value))
	}
	return idfMap
}

func tfidf(documents []string) (map[string]map[string]float64) {
	tfidfMap := make(map[string]map[string]float64)
	idf := idf(documents)

	for _, document := range documents {
		m := make(map[string]float64)
		tf := tf(document)
		for word, docFreq := range idf {
			m[word] = float64(tf[word]) * docFreq
		}
		tfidfMap[document] = m
	}

	return tfidfMap
}
func tfidfMap(documents map[string]string) (map[string]map[string]float64) {
	tfidfMap := make(map[string]map[string]float64)
	var d []string
	for _, document := range documents {
		d = append(d, document)
	}

	idf := idf(d)
	var count int
	for key, document := range documents {
		m := make(map[string]float64)
		tf := tf(document)
		for word, docFreq := range idf {
			m[word] = float64(tf[word]) * docFreq
		}
		tfidfMap[key] = m
		count++
		fmt.Printf("\r%d of %d",count , len(documents) )
	}

	return tfidfMap
}
