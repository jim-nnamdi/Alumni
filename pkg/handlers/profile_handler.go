package handlers

import (
	"net/http"

	"github.com/jim-nnamdi/jinx/pkg/database/mysql"
	"github.com/jim-nnamdi/jinx/pkg/model"
	"go.uber.org/zap"
)

var _ http.Handler = &profileHandler{}
var profileTTL = 3600

type profileHandler struct {
	logger      *zap.Logger
	mysqlclient mysql.Database
}

func NewProfileHandler(logger *zap.Logger, mysqlclient mysql.Database) *profileHandler {
	return &profileHandler{
		logger:      logger,
		mysqlclient: mysqlclient,
	}
}

func (handler *profileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	profileres := map[string]interface{}{}
	userInfo, ok := model.FromContext(r.Context())
	if !ok {
		profileres["err"] = "please sign in to access this page"
		handler.logger.Debug("unauthorized user")
		w.Write(GetSuccessResponse(profileres["err"], profileTTL))
		return
	}
	profileres["id"] = userInfo.Id
	profileres["username"] = userInfo.Username
	profileres["email"] = userInfo.Email
	profileres["phone"] = userInfo.Phone
	profileres["session_key"] = userInfo.SessionKey
	profileres["wallet_balance"] = userInfo.WalletBalance
	w.Write(GetSuccessResponse(profileres, profileTTL))
}
