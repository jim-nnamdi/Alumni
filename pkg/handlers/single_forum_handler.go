package handlers

import (
	"log"
	"net/http"

	"github.com/jim-nnamdi/jinx/pkg/database/mysql"
)

var _ http.Handler = &sforumStruct{}

type sforumStruct struct {
	Log *log.Logger
	Db  mysql.Database
}

func NewSForumStruct(log *log.Logger, Db mysql.Database) *sforumStruct {
	return &sforumStruct{
		Log: log,
		Db:  Db,
	}
}

func (fs *sforumStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		slug = r.FormValue("slug")
		sfp  = map[string]string{}
	)
	if slug == "" {
		fs.Log.Printf("'%s'\n", "slug is empty")
		return
	}

	get_single_forum_post, err := fs.Db.GetSingleForumPost(r.Context(), slug)
	if err != nil {
		fs.Log.Printf("'%s'\n", err)
		sfp["error"] = err.Error()
		w.Write(GetErrorResponseBytes(sfp, 30, err))
		return
	}

	if get_single_forum_post != nil {
		forum_resp := map[string]interface{}{}
		forum_resp["id"] = get_single_forum_post.Id
		forum_resp["title"] = get_single_forum_post.Title
		forum_resp["description"] = get_single_forum_post.Description
		forum_resp["slug"] = get_single_forum_post.Slug
		forum_resp["created_at"] = get_single_forum_post.CreatedAt
		forum_resp["updated_at"] = get_single_forum_post.UpdatedAt
		w.Write(GetSuccessResponse(forum_resp, 30))
	} else {
		sfp["error"] = "no post data"
		w.Write(GetSuccessResponse(sfp, 30))
	}

}
