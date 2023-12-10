package middleware

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/jim-nnamdi/jinx/pkg/database/mysql"
	"github.com/jim-nnamdi/jinx/pkg/model"
	"go.uber.org/zap"
)

type SessionMiddleware struct {
	logger      *zap.Logger
	mysqlclient mysql.Database
}

func NewSessionMiddleware(logger *zap.Logger, mysqlclient mysql.Database) *SessionMiddleware {
	return &SessionMiddleware{
		logger:      logger,
		mysqlclient: mysqlclient,
	}
}

type response struct {
	Data    interface{} `json:"data,omitempty"`
	Err     string      `json:"err,omitempty"`
	Success bool        `json:"success"`
	TTL     int         `json:"ttl"`
}

func GetSuccessResponse(data interface{}, ttl int) []byte {
	resp := &response{
		Success: true,
		TTL:     ttl,
	}

	if data == nil || (reflect.ValueOf(data).Kind() == reflect.Ptr && reflect.ValueOf(data).IsNil()) {
		resp.Data = nil
	} else {
		resp.Data = data
	}
	responseBytes, _ := json.Marshal(resp)
	return responseBytes
}

func (smw *SessionMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessionkey = r.FormValue("session_key")
		var sessionresp = map[string]interface{}{}
		uInfo := &model.User{}
		getBySession, err := smw.mysqlclient.GetBySessionKey(r.Context(), sessionkey)
		if err != nil {
			sessionresp["err"] = "no session key passed"
			smw.logger.Debug("cannot fetch user by sessionkey", zap.Any("fetch user error", err))
			w.Write(GetSuccessResponse(sessionresp["err"], 30))
			return
		}
		uInfo.Username = getBySession.Username
		uInfo.Country = getBySession.Country
		uInfo.Email = getBySession.Email
		uInfo.Phone = getBySession.Phone
		uInfo.SessionKey = getBySession.SessionKey
		uInfo.Id = getBySession.Id
		uInfo.WalletBalance = getBySession.WalletBalance
		ctx := model.NewContext(r.Context(), uInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
