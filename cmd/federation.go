package main

import (
	"os"
	"sublinks/federation/internal/http"
)

func main() {
	http.RunServer()
	os.Exit(0)
}
