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

func validateEmail(email string) (bool, error) {
	if len(email) > 20 {
		return false, fmt.Errorf("email exceeds required length")
	}
	parts := strings.Split(strings.ToLower(email), "@")
	if len(parts) != 2 {
		return false, fmt.Errorf("email must contain @ symbol")
	}
	local, domain := parts[0], parts[1]
	if len(local) == 0 ||
		len(domain) == 0 {
		return false, fmt.Errorf("local or domain cannot be empty")
	}
	prev_char := rune(0)
	for _, char := range local {
		if strings.ContainsRune("!#$%&'*+-/=?^_`{|}~.", char) {
			if char == prev_char && char != '-' {
				return false, fmt.Errorf("cannot contain special chars before domain")
			}
		}
		prev_char = char
	}
	if strings.ContainsAny(email, " ") {
		return false, fmt.Errorf("email cannot contain spaces")
	}
	if len(local) > 64 || len(domain) > 255 {
		return false, fmt.Errorf("local part or domain part length exceeds the limit in the email")
	}
	return false, nil
}

func (handler *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		username        = r.FormValue("username")
		password        = r.FormValue("password")
		email           = r.FormValue("email")
		gradyear        = r.FormValue("grad_year")
		phone           = r.FormValue("phone")
		degree          = r.FormValue("degree")
		currentjob      = r.FormValue("current_job")
		linkedinprofile = r.FormValue("linkedin_profile")
		twitterprofile  = r.FormValue("twitter_profile")
		dataresp        = map[string]interface{}{}
	)

	// there should be a frontend validation for all fields
	// the backend would assist to catch empty fields if the
	// frontend validation is compromised.
	if username == "" || password == "" || degree == "" || phone == "" || email == "" {
		handler.logger.Error("some fields are empty")
		dataresp["err"] = "some fields are empty"
		w.Write(GetSuccessResponse(dataresp, registerTTL))
		return
	}

	// ensure the password is greater than 7 values
	// also ensure that it has special characters
	specialchars := strings.ContainsAny(password, "$ % @")
	passwdcount := len(password)
	if !specialchars {
		handler.logger.Error("password must contain special characters")
		dataresp["err"] = "password must contain special characters"
		w.Write(GetSuccessResponse(dataresp, registerTTL))
		return
	}
	if passwdcount <= 7 {
		handler.logger.Error("password must contain at least 8 characters")
		dataresp["err"] = "password must contain at least 8 characters"
		w.Write(GetSuccessResponse(dataresp, registerTTL))
		return
	}

	// we need to hash the password to avoid security issues
	// and then re-hash it when the user wants to login
	hashed_password, err := hashPassword(password)
	if err != nil {
		handler.logger.Error("cannot hash password", zap.String("hashed password error", err.Error()))
		fmt.Printf("cannot hash password(%s)", password)
		return
	}
	newsessionkey := createSessionKey(email, time.Now())
	sanitize_email, err := validateEmail(email)
	if err != nil || !sanitize_email {
		handler.logger.Error("email was malformed!", zap.Error(err))
		return
	}
	createUser, err := handler.mysqlclient.CreateUser(r.Context(), username, hashed_password, email, degree, gradyear, currentjob, phone, newsessionkey, "", linkedinprofile, twitterprofile)
	if err != nil || !createUser {
		dataresp["err"] = "cannot register user, try again"
		handler.logger.Error("could not create user", zap.Any("error", err))
		w.Write(GetSuccessResponse(dataresp, registerTTL))
		return
	}
	dataresp["username"] = username
	dataresp["email"] = email
	dataresp["degree"] = degree
	dataresp["phone"] = phone
	dataresp["session_key"] = newsessionkey
	handler.logger.Error("user successfully created", zap.Bool("registration success", createUser))
	w.Write(GetSuccessResponse(dataresp, registerTTL))
}
