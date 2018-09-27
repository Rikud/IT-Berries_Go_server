package gameServer

import (
	"IT-Berries/controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type GameServer struct {
	router *mux.Router
	server *http.Server
}

func (s *GameServer) bindHandlers() {
	for path, handle := range controllers.Handlers{
		s.router.Handle(path, handle)
	}
}

func (s *GameServer) Prepare () {
	s.router = mux.NewRouter()
	s.bindHandlers()
}

func (s *GameServer) Start() {
	s.server = &http.Server{
		Handler:      s.router,
		Addr:         "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(s.server.ListenAndServe())
}