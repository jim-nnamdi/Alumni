package handlers

import (
	"net/http"

	"github.com/jim-nnamdi/jinx/pkg/database/mysql"
	"go.uber.org/zap"
)

var _ http.Handler = &loginHandler{}
var loginTTL = 30

type loginHandler struct {
	logger      *zap.Logger
	mysqlclient mysql.Database
}

func NewLoginHandler(logger *zap.Logger, mysqlclient mysql.Database) *loginHandler {
	return &loginHandler{
		logger:      logger,
		mysqlclient: mysqlclient,
	}
}

func (handler *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		email    = r.FormValue("email")
		password = r.FormValue("password")
		loginres = map[string]interface{}{}
	)

	checkuser, err := handler.mysqlclient.GetUserByEmail(r.Context(), email)
	if err != nil {
		loginres["err"] = "user does not exist"
		handler.logger.Error("user does not exist", zap.Any("checkuser", err))
		w.Write(GetSuccessResponse(loginres["err"], loginTTL))
		return
	}
	if checkuser != nil {
		if checkuser.Id > 0 {
			handler.logger.Debug("found user", zap.Bool("user found", true))
			_ = CheckPasswordHash(password, checkuser.Password)
			loginnow, err := handler.mysqlclient.CheckUser(r.Context(), email, checkuser.Password)
			if err != nil {
				loginres["err"] = "email or password incorrect"
				handler.logger.Error("email or password incorrect", zap.Any("login response", "email or password incorrect"))
				w.Write(GetSuccessResponse(loginres["err"], loginTTL))
				return
			}
			if loginnow != nil {
				loginres["username"] = loginnow.Username
				loginres["email"] = loginnow.Email
				loginres["phone"] = loginnow.Phone
				loginres["degree"] = loginnow.Degree
				loginres["grad_year"] = loginnow.GradYear
				loginres["current_job"] = loginnow.CurrentJob
				loginres["profile_picture"] = loginnow.ProfilePicture
				loginres["linkedin_profile"] = loginnow.LinkedinProfile
				loginres["twitter_profile"] = loginnow.TwitterProfile
				w.Write(GetSuccessResponse(loginres, loginTTL))
			}
		}
	}
}
