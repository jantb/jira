package main

import (
	"strings"
)

func toLower(words string) []strings {
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
	for _, stopword := range stopwords {
		for _, word := range words {
			if stopword == word {
				continue
			}
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

func set(words []string) []string {
	set := make([]string, len(words))
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

func tf(string string) (map[string]int) {
	words := removeNorwegianStopwords(toLower(string))
	return wf(words)
}
func idf(string string) (map[string]int) {
	words := set(removeNorwegianStopwords(toLower(string)))
	return wf(words)
}
