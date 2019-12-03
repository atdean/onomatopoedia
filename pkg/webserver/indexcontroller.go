package webserver

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/atdean/onomatopoedia/pkg/repositories"
)

// IndexController contains route handlers for the index page
type IndexController struct {
	SqlPool *sql.DB
}

// GetIndexHandler servers a GET request for route /
func (ctrl *IndexController) GetIndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	entriesRepo := repositories.EntryRepository{SqlPool: ctrl.SqlPool}
	entry, err := entriesRepo.GetByID(1)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("Entry ID: %d, UserID: %d, Slug: %s, DisplayName: %s\n",
			entry.ID, entry.UserID, entry.Slug, entry.DisplayName)
	}

	renderer := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/index.html",
	))
	err = renderer.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err)
	}
}
