package main

import (
	"errors"
	"fmt"
)

type Fact struct {
	u1     string
	cValue float64
	u2     string
}

type Query struct {
	value float64
	from  string
	to    string
}

type Graph [][]float64

func main() {
	facts := []Fact{
		{"mtr", 3.28, "ft"},
		{"ft", 12, "in"},
		{"hr", 60, "min"},
		{"min", 60, "sec"},
	}
	factMap := map[string]int{
		"mtr": 0,
		"ft":  1,
		"in":  2,
		"hr":  3,
		"min": 4,
		"sec": 5,
	}
	g := createGraph(len(factMap))
	pg := populateGraph(g, facts, factMap)
	fmt.Printf("%v", pg)
}

func parseFacts(facts []Fact) (factMap map[string]int) {
	factMap = make(map[string]int)
	var vNum int

	for _, fact := range facts {

		if _, exists := factMap[fact.u1]; !exists {
			factMap[fact.u1] = vNum
			vNum += 1
		}

		if _, exists := factMap[fact.u2]; !exists {
			factMap[fact.u2] = vNum
			vNum += 1
		}
	}

	return
}

func sliceContains(haystack []string, needle string) bool {
	for _, hay := range haystack {
		if hay == needle {
			return true
		}
	}

	return false
}

func createGraph(size int) Graph {
	outer := make([][]float64, size)
	for i := range outer {
		outer[i] = make([]float64, size)
	}

	return outer
}

func populateGraph(graph Graph, facts []Fact, factMap map[string]int) Graph {
	for _, fact := range facts {
		fromUnit := factMap[fact.u1]
		toUnit := factMap[fact.u2]

		graph[fromUnit][toUnit] = fact.cValue     // one direction
		graph[toUnit][fromUnit] = 1 / fact.cValue // opposite direction
	}

	return graph
}

func convert(graph Graph, factMap map[string]int, q Query) (float64, error) {
	visited := createGraph(len(graph))
	frmV, toV := factMap[q.from], factMap[q.to]
	m, err := findPathMultiplier(graph, visited, frmV, toV)
	return q.value * m, err
}

func findPathMultiplier(g, visited Graph, frmV, toV int) (float64, error) {
	// base case
	if g[frmV][toV] != 0 {
		return g[frmV][toV], nil
	}

	for ci := 0; ci < len(g); ci++ {
		if visited[frmV][ci] == 0 && g[frmV][ci] != 0 { // unvisited & non-empty edge
			visited[frmV][ci] = 1
			m, err := findPathMultiplier(g, visited, ci, toV)
			if err == nil {
				return g[frmV][ci] * m, nil
			}
		}
	}
	return 0, errors.New("not convertible")
}
