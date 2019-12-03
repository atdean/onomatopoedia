package webserver

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// LoginController handles all routes for logging in and out of the application
type LoginController struct {
}

// GetLoginHandler displays the "log in" page
func (ctrl *LoginController) GetLoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderer := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/auth/login.html",
	))
	renderer.ExecuteTemplate(w, "base", nil)
}
