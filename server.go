package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/wolfeidau/proxy-auth/assets"
)

const (
	// PathPrefix the path used by this authentication module
	PathPrefix = "/auth"
	// LoginURL is the URL which users are redirected to if they aren't logged in
	LoginURL = "/auth/login"
)

// Server the authentication server
type Server struct {
	mux          *mux.Router
	sessionStore sessions.Store
	oauthConfig  *oauth2.Config
}

// NewServer creates a server with the standard endpionts registered
func NewServer(sessionStore sessions.Store) *Server {

	conf := &oauth2.Config{
		ClientID:     defaultGitHubConfig.ClientID,
		ClientSecret: defaultGitHubConfig.ClientSecret,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	r := mux.NewRouter().PathPrefix("/auth").Subrouter()
	s := &Server{mux: r, sessionStore: sessionStore, oauthConfig: conf}

	r.HandleFunc("/github/authorize", s.loginGitHubAuthorise).Methods("GET")
	r.HandleFunc("/login", s.loginHandler).Methods("GET")
	r.HandleFunc("/logout", s.logoutHandler).Methods("GET")
	r.HandleFunc("/github/redirect", s.redirectGitHubHandler).Methods("GET")
	r.HandleFunc("/{asset}", s.assetHandler).Methods("GET")

	return s
}

// GetMux returns the mux with the standard http handlers already registered
func (s *Server) GetMux() *mux.Router {
	return s.mux
}

// CheckSession middleware function to validate the session cookie is set
func CheckSession(handler http.Handler, store sessions.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth")

		fmt.Printf("CheckSession email=%v url=%s\n", session.Values["email"], r.URL.String())

		if !strings.HasPrefix(r.URL.String(), "/auth/") && session.Values["email"] == nil {

			http.Redirect(w, r, LoginURL, http.StatusFound)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (s *Server) loginGitHubAuthorise(w http.ResponseWriter, r *http.Request) {
	fmt.Println("loginGitHubAuthorise")
	session, _ := s.sessionStore.Get(r, "auth")
	state, _ := generateState()

	// assign the state variable in the session
	session.Values["state"] = state

	session.Save(r, w)

	url := s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)

	http.Redirect(w, r, url, http.StatusFound)
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("loginHandler")
	//session, _ := s.sessionStore.Get(r, "auth")

	buf, err := assets.Asset("index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(buf)
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logoutHandler")
	session, _ := s.sessionStore.Get(r, "auth")

	delete(session.Values, "name")
	delete(session.Values, "email")

	session.Save(r, w)

	http.Redirect(w, r, LoginURL, http.StatusFound)

}

type redirectOptions struct {
	Code  string `schema:"code"`
	State string `schema:"state"`
}

func (s *Server) redirectGitHubHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("redirectGitHubHandler")
	session, _ := s.sessionStore.Get(r, "auth")
	state := session.Values["state"]

	decoder := schema.NewDecoder()
	reqOpts := new(redirectOptions)
	decoder.Decode(reqOpts, r.URL.Query())

	// check state
	if state != reqOpts.State {
		//fmt.Printf("woops got %s expected %s\n", state, reqOpts.State)
		err := fmt.Errorf("state missing")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		delete(session.Values, "state")
	}

	tok, err := s.oauthConfig.Exchange(oauth2.NoContext, reqOpts.Code)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := s.oauthConfig.Client(oauth2.NoContext, tok)

	resp, err := client.Get("https://api.github.com/user")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := DecodeGitHubUser(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["name"] = user.Name
	session.Values["email"] = user.Email

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) assetHandler(w http.ResponseWriter, r *http.Request) {
	assetName := mux.Vars(r)["asset"]
	fmt.Printf("url=%s asset=%s\n", r.URL.String(), assetName)

	buf, err := assets.Asset(assetName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Write(buf)
}
