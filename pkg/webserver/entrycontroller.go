package webserver

import (
	"github.com/atdean/onomatopoedia/pkg/repositories"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
)

type EntryController struct {
	App *App
}

func newEntryController(app *App) *EntryController {
	return &EntryController{
		App: app,
	}
}

func (ctrl *EntryController) GetSingleEntryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	entriesRepo := repositories.NewEntryRepository(ctrl.App.SqlPool)
	entry, err := entriesRepo.GetBySlug(ps.ByName("slug"))
	if err != nil {
		log.Println(err)
	}

	data := map[string]interface{}{
		"PageTitle": entry.DisplayName,
		"Entry":     entry,
	}

	renderer := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/single.html",
	))
	if err = renderer.ExecuteTemplate(w, "base", data); err != nil {
		log.Println(err)
	}
}