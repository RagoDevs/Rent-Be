package main

import "net/http"

func (app *application) test(w http.ResponseWriter , r *http.Request) {
    w.Write([]byte("He we test authentication"))
}
