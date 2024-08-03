package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/jim-nnamdi/jinx/pkg/database/mysql"
	"github.com/jim-nnamdi/jinx/pkg/model"
)

type IChat interface {
	ServeHttp(w http.ResponseWriter, r *http.Request)
}

var _ IChat = &ichatStruct{}

type ichatStruct struct {
	Log log.Logger
	DB  mysql.Database
}

func (cs *ichatStruct) ServeHttp(w http.ResponseWriter, r *http.Request) {
	var (
		recipient = r.FormValue("recv_email")
		message   = r.FormValue("message")
	)
	current_user, ok := model.FromContext(r.Context())
	if !ok {
		cs.Log.Printf("'%s'\n", "not authenticated")
		cs.Log.Printf("'Email: %s'\n", current_user.Email)
		return
	}
	recv_user, err := cs.DB.GetUserByEmail(r.Context(), recipient)
	if err != nil {
		cs.Log.Printf("'%s'\n", err)
		cs.Log.Printf("'%s'\n", "check if recipient exists")
		return
	}
	send_chat, err := cs.DB.SendMessage(r.Context(), current_user.Id, recv_user.Id, message, time.Now(), time.Now())
	if err != nil {
		cs.Log.Printf("'%s'\n", "could not send message to recipient")
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
