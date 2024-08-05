package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jim-nnamdi/jinx/pkg/database/mysql"
	"github.com/jim-nnamdi/jinx/pkg/middleware"
	"go.uber.org/zap"
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
	Logger         *zap.Logger
	Mysqlclient    mysql.Database
}

func (server *GracefulShutdownServer) getRouter() *mux.Router {
	// Apply middleware to specific routes
	router := mux.NewRouter()
	router.Handle("/login", server.LoginHandler).Methods(http.MethodPost)
	router.Handle("/register", server.RegisterHandler).Methods((http.MethodPost))

	// Wrap profile and forum routes with session middleware
	sessionMiddleware := middleware.NewSessionMiddleware(server.Logger, server.Mysqlclient)
	router.Handle("/profile", sessionMiddleware.Middleware(server.ProfileHandler)).Methods(http.MethodGet)
	router.Handle("/forum", sessionMiddleware.Middleware(server.AllForumHandler)).Methods(http.MethodGet)
	router.Handle("/forum-post", sessionMiddleware.Middleware(server.SingleForumHandler)).Methods(http.MethodGet)
	router.Handle("/add-forum-post", sessionMiddleware.Middleware(server.AddForumHandler)).Methods(http.MethodPost)
	router.Handle("/chat", sessionMiddleware.Middleware(server.ChatHandler)).Methods(http.MethodPost)
	router.SkipClean(true)
	return router
}

func (server *GracefulShutdownServer) Start() {
	logger, _ := zap.NewDevelopment()
	router := server.getRouter()
	server.httpServer = &http.Server{
		Addr:         server.HTTPListenAddr,
		WriteTimeout: server.WriteTimeout,
		ReadTimeout:  server.ReadTimeout,
		IdleTimeout:  server.IdleTimeout,
		Handler:      router,
	}
	logger.Sugar().Info(fmt.Sprintf("listening and serving on %s", server.HTTPListenAddr))
	if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("server failed to start", zap.Error(err))
	}
}
