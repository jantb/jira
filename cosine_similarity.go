package main

import (
	"math"
)

func similarityVector(a, b []float64, dot float64) float64 {
	var normA float64
	var normB float64
	for _, av := range a {
		normA += math.Pow(av, 2)
	}
	for _, av := range b {
		normB += math.Pow(av, 2)
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Keys exstracts keys from map
func Keys(m map[string]float64) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Vals exstracts vals from map
func Vals(m map[string]float64) (vals []float64) {
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

func getIntersection(a, b map[string]float64) (keys []string) {
	for _, value := range Keys(a) {
		if _, ok := b[value]; ok {
			keys = append(keys, value)
		}
	}
	return keys
}

func dot(a, b map[string]float64, intersection []string) (dotProduct float64) {
	for _, av := range intersection {
		dotProduct += a[av] * b[av]
	}
	return dotProduct
}

func similarity(a, b map[string]float64) float64 {
	dot := dot(a, b, getIntersection(a, b))
	return similarityVector(Vals(a), Vals(b), dot)
}
