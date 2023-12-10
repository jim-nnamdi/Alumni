package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jim-nnamdi/jinx/pkg/database/mysql"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var _ http.Handler = &registerHandler{}

var registerTTL = 60

type registerHandler struct {
	logger      *zap.Logger
	mysqlclient mysql.Database
}

func NewRegisterHandler(logger *zap.Logger, mysqlclient mysql.Database) *registerHandler {
	return &registerHandler{
		logger:      logger,
		mysqlclient: mysqlclient,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return password, err
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (handler *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		username = r.FormValue("username")
		password = r.FormValue("password")
		email    = r.FormValue("email")
		country  = r.FormValue("country")
		phone    = r.FormValue("phone")
		dataresp = map[string]interface{}{}
	)

	// there should be a frontend validation for all fields
	// the backend would assist to catch empty fields if the
	// frontend validation is compromised.
	if username == "" || password == "" || country == "" || phone == "" || email == "" {
		handler.logger.Debug("some fields are empty")
		dataresp["err"] = "some fields are empty"
		w.Write(GetSuccessResponse(dataresp, registerTTL))
		return
	}

	// ensure the password is greater than 7 values
	// also ensure that it has special characters
	specialchars := strings.ContainsAny(password, "$ % @")
	passwdcount := len(password)
	if !specialchars {
		handler.logger.Debug("password must contain special characters")
		dataresp["err"] = "password must contain special characters"
		w.Write(GetSuccessResponse(dataresp, registerTTL))
		return
	}
	if passwdcount <= 7 {
		handler.logger.Debug("password must contain at least 8 characters")
		dataresp["err"] = "password must contain at least 8 characters"
		w.Write(GetSuccessResponse(dataresp, registerTTL))
		return
	}

	// we need to hash the password to avoid security issues
	// and then re-hash it when the user wants to login
	hashed_password, err := hashPassword(password)
	if err != nil {
		handler.logger.Debug("cannot hash password", zap.String("hashed password error", err.Error()))
		fmt.Printf("cannot hash password(%s)", password)
		return
	}
	newsessionkey := createSessionKey(email, time.Now())
	createUser, err := handler.mysqlclient.CreateUser(r.Context(), username, hashed_password, email, country, phone, newsessionkey, 0.0)
	if err != nil || !createUser {
		dataresp["err"] = "cannot register user, try again"
		handler.logger.Debug("could not create user", zap.Any("error", err))
		w.Write(GetSuccessResponse(dataresp, registerTTL))
		return
	}
	dataresp["username"] = username
	dataresp["email"] = email
	dataresp["country"] = country
	dataresp["phone"] = phone
	dataresp["session_key"] = newsessionkey
	handler.logger.Debug("user successfully created", zap.Bool("registration success", createUser))
	w.Write(GetSuccessResponse(dataresp, registerTTL))
}
