package webserver

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/atdean/onomatopoedia/pkg/repositories"
)

type IndexController struct {
	App *App
}

func NewIndexController(app *App) *IndexController {
	return &IndexController{
		App: app,
	}
}

func (ctrl *IndexController) GetIndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	entriesRepo := repositories.NewEntryRepository(ctrl.App.SqlPool)
	entries, err := entriesRepo.GetMostRecent(10, 1)
	if err != nil {
		log.Println(err)
	}

	data := map[string]interface{}{
		"PageTitle": "Home",
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
