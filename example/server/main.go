package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/wolfeidau/proxy-auth"
)

func main() {
	// setup a store, in our case one using secure cookies
	store := sessions.NewCookieStore([]byte("something-very-secret"))
	s := auth.NewServer(store)

	// configure a mux
	r := mux.NewRouter()
	r.PathPrefix(auth.PathPrefix).Handler(s.GetMux())

	// add a wrapper to check the session for each request
	o := auth.CheckSession(r, store)

	// listen to the network
	http.ListenAndServe(":5000", o)
}
