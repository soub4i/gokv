package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

type Store struct {
	store sync.Map
}

type APP struct {
	router *mux.Router
	KVS    *Store
}

func (app *APP) setup() {
	app.router.HandleFunc("/health", ping)
	app.router.HandleFunc("/set/{key}/{value}", app.SetHandler)
	app.router.HandleFunc("/get/{key}", app.GETHandler)
	app.router.HandleFunc("/delete/{key}", app.DeleteHandler)
}
func (app *APP) run() {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.setup()
	log.Println("Server is running on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, app.router))

}

func (s *Store) Set(key string, value string) {
	s.store.LoadOrStore(key, value)
}

func (s *Store) Get(key string) (string, bool) {
	value, ok := s.store.Load(key)
	if !ok {
		return "", false
	}
	return value.(string), true
}

func (s *Store) Delete(key string) {
	s.store.Delete(key)
}

func (s *Store) setupRoute(key string) {
	s.store.Delete(key)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func (a *APP) GETHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("key is required"))
	}

	value, found := a.KVS.Get(key)

	if !found {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("value not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}
func (a *APP) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("key is required"))
	}

	a.KVS.Delete(key)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))

}

func (a *APP) SetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value := vars["value"]
	if key == "" || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("key/value is required"))
		return
	}

	a.KVS.Set(key, value)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func main() {

	var app = &APP{
		router: mux.NewRouter(),
		KVS: &Store{
			store: sync.Map{},
		},
	}

	app.run()

}
