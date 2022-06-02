package main

import (
	"testing"
)

func teststripSuffixPrefix(t *testing.T) {

	testSet := []string{"[testResult]", "[testResult", "testResult]", "testResult"}

	tesResult := "testResult"

	for _, test := range testSet {
		if stripSuffixPrefix(test) != tesResult {
			t.Error("For the set: ", test, " was expect: ", tesResult, " got: ", stripSuffixPrefix(test))
		}
	}
}
