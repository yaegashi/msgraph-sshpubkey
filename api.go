package main

import (
	"fmt"
	"html"
	"net/http"
)

func (app *ServeCmd) HandleAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
