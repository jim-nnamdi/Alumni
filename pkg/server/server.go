package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type GracefulShutdownServer struct {
	HTTPListenAddr  string
	RegisterHandler http.Handler // register
	LoginHandler    http.Handler // login
	ProfileHandler  http.Handler // profile
	HomeHandler     http.Handler

	AddForumHandler    http.Handler // add forum post
	AllForumHandler    http.Handler // get all posts
	SingleForumHandler http.Handler // get one post
	ChatHandler        http.Handler // chat a user

	httpServer     *http.Server
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	IdleTimeout    time.Duration
	HandlerTimeout time.Duration
}

func (server *GracefulShutdownServer) getRouter() *mux.Router {
	router := mux.NewRouter()
	router.Handle("/login", server.LoginHandler)
	router.Handle("/register", server.RegisterHandler)
	router.Handle("/profile", server.ProfileHandler)
	router.Handle("/forum", server.AllForumHandler)
	router.Handle("/forum-post", server.SingleForumHandler)
	router.Handle("/add-forum-post", server.AddForumHandler)
	router.Handle("/chat", server.ChatHandler)
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
