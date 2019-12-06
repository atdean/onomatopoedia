package webserver

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/atdean/onomatopoedia/pkg/repositories"
)

type IndexController struct {
	SqlPool *sql.DB
}

func NewIndexController(db *sql.DB) *IndexController {
	return &IndexController{
		SqlPool: db,
	}
}

func (ctrl *IndexController) GetIndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	entriesRepo := repositories.EntryRepository{SqlPool: ctrl.SqlPool}
	entries, err := entriesRepo.GetMostRecent(10, 1)
	if err != nil {
		log.Println(err)
	}

	data := map[string]interface{}{
		"Entries": entries,
	}

	renderer := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/index.html",
	))
	if err = renderer.ExecuteTemplate(w, "base", data); err != nil {
		log.Println(err)
	}
}
