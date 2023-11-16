package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type GracefulShutdownServer struct {
	HTTPListenAddr  string
	RegisterHandler http.Handler
	LoginHandler    http.Handler
	ProfileHandler  http.Handler
	HomeHandler     http.Handler

	httpServer     *http.Server
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	IdleTimeout    time.Duration
	HandlerTimeout time.Duration
}

func (server *GracefulShutdownServer) getRouter() *mux.Router {
	router := mux.NewRouter()

	router.SkipClean(true)
	return router
}

func (server *GracefulShutdownServer) Start() {
	router := server.getRouter()
	server.httpServer = &http.Server{
		Addr:         server.HTTPListenAddr,
		WriteTimeout: server.WriteTimeout,
		ReadTimeout:  server.ReadTimeout,
		IdleTimeout:  server.IdleTimeout,
		Handler:      router,
	}
	server.httpServer.ListenAndServe()
}
