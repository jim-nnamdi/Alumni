package handlers

import (
	"log"
	"net/http"

	"github.com/jim-nnamdi/jinx/pkg/database/mysql"
)

type SForum interface {
	ServeHttp(w http.ResponseWriter, r *http.Request)
}

var _ SForum = &sforumStruct{}

type sforumStruct struct {
	Log log.Logger
	Db  mysql.Database
}

func (fs *sforumStruct) ServeHttp(w http.ResponseWriter, r *http.Request) {
	var (
		slug = r.FormValue("slug")
	)
	if slug == "" {
		fs.Log.Printf("'%s'\n", "slug is empty")
		return
	}

	get_single_forum_post, err := fs.Db.GetSingleForumPost(r.Context(), slug)
	if err != nil {
		fs.Log.Printf("'%s'\n", err)
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
		w.Write(GetSuccessResponse([]struct{}{}, 30))
	}

}
