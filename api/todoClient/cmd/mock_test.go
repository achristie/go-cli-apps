package cmd

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
)

var (
	//go:embed testdata/ResultsMany.json
	ResultsMany string

	//go:embed testdata/ResultsOne.json
	ResultsOne string

	//go:embed testdata/NoResults.json
	NoResults string
)

var testResponses = map[string]struct {
	Status int
	Body   string
}{
	"resultsMany": {
		Status: http.StatusOK,
		Body:   ResultsMany,
	},
	"resultsOne": {
		Status: http.StatusOK,
		Body:   ResultsOne,
	},
	"noResults": {
		Status: http.StatusOK,
		Body:   NoResults,
	},
	"root": {
		Status: http.StatusOK,
		Body:   "API HERE",
	},
	"notFound": {
		Status: http.StatusNotFound,
		Body:   "404 - not found",
	},
}

func mockServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)

	return ts.URL, func() { ts.Close() }
}
