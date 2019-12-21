package webserver

import (
	"database/sql"
	"fmt"
	"github.com/atdean/onomatopoedia/pkg/models"
	"github.com/gomodule/redigo/redis"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"

	"github.com/atdean/onomatopoedia/pkg/repositories"
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
func (ctrl *AuthController) GetSignupHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderer := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/auth/signup.html",
	))
	if err := renderer.ExecuteTemplate(w, "base", nil); err != nil {
		ctrl.App.ServerError(w, err)
		return
	}
}

// POST /signup
func (ctrl *AuthController) PostSignupHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO :: Sanitize and validate form data
	username := r.PostFormValue("username")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	repo := repositories.NewUserRepository(ctrl.App.SqlPool)
	newUser, err := repo.CreateNewUser(username, email, password)
	if err != nil {
		ctrl.App.ServerError(w, err)
		return
	}

	// TODO :: Log the new user in and redirect them to their profile page
	fmt.Fprintf(w, "New user created!\n%v\n", newUser)
}

// GET /login
func (ctrl *AuthController) GetLoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderer := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/auth/login.html",
	))
	if err := renderer.ExecuteTemplate(w, "base", nil); err != nil {
		ctrl.App.ServerError(w, err)
		return
	}
}

// POST /login
func (ctrl *AuthController) PostLoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	if len(username) <= 0 || len(password) <= 0 {
		//ctrl.App.ClientError(w, http.StatusUnauthorized)
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	} else {
		repo := repositories.NewUserRepository(ctrl.App.SqlPool)

		user, err := repo.GetByUsername(username)
		if err != nil || user == nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("<div>Username not found.</div>"))
			} else {
				ctrl.App.ServerError(w, err)
				return
			}
		} else {
			if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("<div>Incorrect password.</div>"))
			} else {
				// Create a random-number session token
				sessionToken := uuid.NewV4().String()
				fmt.Println(sessionToken)

				_, err := ctrl.App.RedisConn.Do("SETEX", sessionToken, "21600", username)
				if err != nil {
					ctrl.App.ServerError(w, err)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name: "session_token",
					Value: sessionToken,
					// TODO :: take the session expiration time from config/env var
					Expires: time.Now().Add(6 * time.Hour),
				})

				http.Redirect(w, r, "/", http.StatusOK)
			}
		}
	}
}

func (ctrl *AuthController) GetLoggedInUser(r *http.Request) (*models.User, error) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil && err == http.ErrNoCookie {
		// No user is logged on
		return nil, err
	}

	sessionToken := sessionCookie.Value

	username, err := redis.String(ctrl.App.RedisConn.Do("GET", sessionToken))
	if err != nil {
		ctrl.App.ErrorLogger.Printf("Unable to fetch session token from redis for %s\n", sessionToken)
		return nil, err
	} else if username == "" {
		ctrl.App.ErrorLogger.Printf("Username fetched from session token is blank for %s\n", sessionToken)
		return nil, fmt.Errorf("Could not validate logged in user.")
	}

	repo := repositories.NewUserRepository(ctrl.App.SqlPool)
	user, err := repo.GetByUsername(username)
	if err != nil {
		ctrl.App.ErrorLogger.Printf("AuthController->GetLoggedInUser: Username not found in database.")
		return nil, err
	}

	return user, nil
}