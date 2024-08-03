package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jim-nnamdi/jinx/pkg/database/mysql"
)

var _ http.Handler = &forumStruct{}

type forumStruct struct {
	Log *log.Logger
	Db  mysql.Database
}

func NewForumStruct(log *log.Logger, Db mysql.Database) *forumStruct {
	return &forumStruct{
		Log: log,
		Db:  Db,
	}
}

func (fs *forumStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		title       = r.FormValue("title")
		description = r.FormValue("description")
		author      = r.FormValue("author")
		afp         = map[string]string{}
	)
	if title == "" || description == "" || author == "" {
		fs.Log.Printf("'%s'\n", "title | description | author is empty")
		afp["error"] = "empty title or description or author"
		w.Write(GetErrorResponseBytes(afp, 30, fmt.Errorf("'%s'", "empty title description or author")))
		return
	}

	if len(title) < 5 {
		fs.Log.Printf("'%s'\n", "invalid title length")
		return
	}

	if len(description) > 200 {
		fs.Log.Printf("'%s'\n", "max length of description exceeded")
		return
	}

	slug := strings.Split(title, " ")
	_slug := strings.Join(slug, "")

	add_new_forum_post, err := fs.Db.AddNewForumPost(r.Context(), title, description, author, _slug, time.Now(), time.Now())
	if err != nil {
		fs.Log.Printf("'%s'\n", err)
		afp["error"] = err.Error()
		w.Write(GetErrorResponseBytes(afp, 30, err))
		return
	}

	if add_new_forum_post {
		new_forum_response := map[string]interface{}{}
		new_forum_response["title"] = title
		new_forum_response["author"] = author
		new_forum_response["message"] = "forum post added successfully"
		w.Write(GetSuccessResponse(new_forum_response, 30))
	}
}
