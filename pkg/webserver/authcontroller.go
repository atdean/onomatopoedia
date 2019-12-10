package webserver

import (
	"database/sql"
	"github.com/atdean/onomatopoedia/pkg/repositories"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type AuthController struct {
	App *App
}

func NewAuthController(app *App) *AuthController {
	return &AuthController{
		App: app,
	}
}

// GET /signup

// POST /signup
// Code to hash password, for later use
// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)

// GET /login
func (ctrl *AuthController) GetLoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderer := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/auth/login.html",
	))
	if err := renderer.ExecuteTemplate(w, "base", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<div>Sorry, an internal server error occurred.</div>"))
	}
}

// POST /login
func (ctrl *AuthController) PostLoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	if len(username) <= 0 || len(password) <= 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<div>Username or password not provided.</div>"))
	} else {

		repo := repositories.NewUserRepository(ctrl.App.SqlPool)
		user, err := repo.GetByUsername(username)
		if err != nil || user == nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("<div>Username not found.</div>"))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("<div>Sorry, an internal server error occurred.</div>"))
			}
		} else {
			if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("<div>Incorrect password."))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("<div>You made it!</div>"))
			}
		}
	}
}