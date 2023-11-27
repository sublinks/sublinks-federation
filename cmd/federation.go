package main

import (
	"os"
	"participating-online/sublinks-federation/internal/http"
)

func main() {
	http.RunServer()
	os.Exit(0)
}
