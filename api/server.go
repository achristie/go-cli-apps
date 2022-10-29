package main

import "net/http"

func newMux(todoFile string) http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/", rootHandler)

	return m
}
