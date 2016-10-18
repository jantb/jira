package main

import (
	"math"
)

func similarityVector(a,b []float64)float64{
	var dotProduct float64
	var normA float64
	var normB float64
	for i, av := range a {
		dotProduct += av * b[i]
		normA += math.Pow(av, 2)
		normB += math.Pow(b[i], 2)
	}
	return dotProduct/(math.Sqrt(normA) *math.Sqrt(normB))
}
func similarity(a,b map[string]float64)float64{
	m := make(map[string][]float64)
	for key, va := range a {
		vb :=b[key]
		m[key] = []float64{va,vb}
	}
	for key, vb := range b {
		va :=a[key]
		m[key] = []float64{va,vb}
	}
	var vecA, vecB []float64
	for _, value := range m {
		vecA = append(vecA,value[0])
		vecB = append(vecB,value[1])
	}
	return similarityVector(vecA,vecB)
}
