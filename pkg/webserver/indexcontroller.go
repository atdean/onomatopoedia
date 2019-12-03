package webserver

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// IndexController contains route handlers for the index page
type IndexController struct {
	DB *sql.DB
}

// GetIndexHandler servers a GET request for route /
func (ctrl *IndexController) GetIndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	renderer := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/index.html",
	))
	renderer.ExecuteTemplate(w, "base", nil)
}
