package main

import (
	"errors"
	"math"
	"reflect"
	"testing"
)

type ExpectedOutput struct {
	value float64
	err   error
}

var ErrNotConvertible = errors.New("not convertible")

func TestParseFacts(t *testing.T) {
	facts := []Fact{
		{"mtr", 3.28, "ft"},
		{"ft", 12, "in"},
		{"hr", 60, "min"},
		{"min", 60, "sec"},
	}

	// i wish for a mapping from unit name to vertex number, Mr. Genie!
	want := map[string]int{
		"mtr": 0,
		"ft":  1,
		"in":  2,
		"hr":  3,
		"min": 4,
		"sec": 5,
	}
	got := parseFacts(facts)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestSliceContains(t *testing.T) {
	given := []string{"a", "b", "c", "d"}
	t.Run("slice contains", func(t *testing.T) {
		got := sliceContains(given, "b")
		if got != true {
			t.Error("Got false, want true")
		}
	})
	t.Run("slice does not contains", func(t *testing.T) {
		got := sliceContains(given, "z")
		if got != false {
			t.Error("Got true, want false")
		}
	})
}

func TestCreateGraph(t *testing.T) {
	given := []string{"mtr", "ft", "in"}

	got := createGraph(len(given))
	want := Graph{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestPopulateGraph(t *testing.T) {
	givenGraph := Graph{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
	givenFacts := []Fact{
		{"mtr", 3.28, "ft"},
		{"ft", 12, "in"},
	}
	givenFactMap := map[string]int{"mtr": 0, "ft": 1, "in": 2}

	want := Graph{
		{0, 3.28, 0},
		{1 / 3.28, 0, 12},
		{0, 1.0 / 12, 0},
	}

	got := populateGraph(givenGraph, givenFacts, givenFactMap)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

// acceptance test
func TestConvert(t *testing.T) {
	// given
	givenFactMap := map[string]int{
		"mtr": 0,
		"ft":  1,
		"in":  2,
		"hr":  3,
		"min": 4,
		"sec": 5,
	}
	givenGraph := Graph{
		{0, 3.28, 0, 0, 0, 0},
		{0.30487, 0, 12, 0, 0, 0},
		{0, 0.083333, 0, 0, 0, 0},
		{0, 0, 0, 0, 60, 0},
		{0, 0, 0, 0.0166, 0, 60},
		{0, 0, 0, 0, 0.01666, 0},
	}

	cases := []struct {
		Name     string
		Input    Query
		Expected ExpectedOutput
	}{
		{
			"mtr to in",
			Query{2, "mtr", "in"},
			ExpectedOutput{78.72, nil},
		},
		{
			"in to mtr",
			Query{13, "in", "mtr"},
			ExpectedOutput{0.330, nil},
		},
		{
			"not convertible",
			Query{13, "in", "hr"},
			ExpectedOutput{0, ErrNotConvertible},
		},
		{
			"hr to mins",
			Query{1.5, "hr", "min"},
			ExpectedOutput{90, nil},
		},
		{
			"min to hr",
			Query{30, "min", "hr"},
			ExpectedOutput{0.5, nil},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got, err := convert(givenGraph, givenFactMap, test.Input)

			if test.Expected.err == nil {
				if err != nil {
					t.Fatal("Wasn't expecting error")
				}
				if floatEqual(got, test.Expected.value) {
					t.Errorf("got %f want %f", got, test.Expected)
				}
			} else {
				if err == nil {
					t.Errorf("was expecting error")
				}
			}

		})
	}
}

func floatEqual(got, want float64) bool {
	eqThreshold := 1e-2
	return math.Abs(got-want) > eqThreshold
}

func TestFindPathMultiplier(t *testing.T) {
	given := Graph{
		{0, 3.28, 0},
		{1 / 3.28, 0, 12},
		{0, 1.0 / 12, 0},
	}
	fromVertex, toVertex := 0, 2
	visited := Graph{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
	got, err := findPathMultiplier(given, visited, fromVertex, toVertex)
	if err != nil {
		t.Fatal("wasn't expecting an error!")
	}
	want := 3.28 * 12

	if got != want {
		t.Errorf("got %f want %f", got, want)
	}
}
