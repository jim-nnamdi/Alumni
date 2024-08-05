package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jim-nnamdi/jinx/pkg/database/mysql"
	"github.com/jim-nnamdi/jinx/pkg/model"
)

var _ http.Handler = &ichatStruct{}

type ichatStruct struct {
	Log *log.Logger
	DB  mysql.Database
}

func NewChat(log *log.Logger, Db mysql.Database) *ichatStruct {
	return &ichatStruct{
		Log: log,
		DB:  Db,
	}
}

func (cs *ichatStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		recipient = r.FormValue("recv_email")
		message   = r.FormValue("message")
	)
	current_user, ok := model.FromContext(r.Context())
	if !ok {
		cs.Log.Printf("'%s'\n", "not authenticated")
		cs.Log.Printf("'Email: %s'\n", current_user.Email)
		model_invalid_sess := map[string]string{}
		model_invalid_sess["error"] = "failed to get user"
		model_invalid_sess["db_error"] = "authentication failed"
		w.Write(GetErrorResponseBytes(model_invalid_sess, 30, fmt.Errorf("'%s'", "auth again")))
		return
	}
	recv_user, err := cs.DB.GetUserByEmail(r.Context(), recipient)
	if err != nil {
		cs.Log.Printf("'%s'\n", err)
		cs.Log.Printf("'%s'\n", "check if recipient exists")
		failed_retrieval := map[string]string{}
		failed_retrieval["error"] = "failed to get user"
		failed_retrieval["db_error"] = err.Error()
		w.Write(GetErrorResponseBytes(failed_retrieval, 30, fmt.Errorf("'%s'", err.Error())))
		return
	}
	msg_valid := len(message)
	if msg_valid > 20 {
		cs.Log.Printf("'%s'\n", "max message threshold")
		cs.Log.Printf("'%s'\n", message)
		msg_resp := map[string]string{}
		msg_resp["error"] = "max message threshold"
		w.Write(GetErrorResponseBytes(msg_resp, 30, fmt.Errorf("'%s'", "error sending message")))
		return
	}
	send_chat, err := cs.DB.SendMessage(r.Context(), current_user.Id, recv_user.Id, message)
	if err != nil {
		cs.Log.Printf("'%s'\n", "could not send message to recipient")
		nilc_resp := map[string]string{}
		nilc_resp["error"] = "error sending message"
		nilc_resp["db_error"] = err.Error()
		w.Write(GetErrorResponseBytes(nilc_resp, 30, err))
		return
	}
	if send_chat {
		chatresp := map[string]interface{}{}
		chatresp["sender"] = current_user.Username
		chatresp["receiver"] = recv_user.Username
		chatresp["message"] = message
		chatresp["created_at"] = time.Now()
		chatresp["updated_at"] = time.Now()
		w.Write(GetSuccessResponse(chatresp, 30))
		return
	}
}
