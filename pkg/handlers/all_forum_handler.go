package handlers

import (
	"log"
	"net/http"

	"github.com/jim-nnamdi/jinx/pkg/database/mysql"
)

type AForum interface {
	ServeHttp(w http.ResponseWriter, r *http.Request)
}

var _ AForum = &aforumStruct{}

type aforumStruct struct {
	Log log.Logger
	Db  mysql.Database
}

func (fs *aforumStruct) ServeHttp(w http.ResponseWriter, r *http.Request) {
	get_all_posts, err := fs.Db.GetAllForums(r.Context())
	if err != nil {
		fs.Log.Printf("'%s'\n", err)
		return
	}

	if get_all_posts != nil {
		w.Write(GetSuccessResponse(get_all_posts, 30))
	} else {
		w.Write(GetSuccessResponse([]struct{}{}, 30))
	}

}
