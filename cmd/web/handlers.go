package main

import (
	"errors"
	"fmt"
	"html/template"

	//"html/template"
	"net/http"
	"strconv"

	"github.com/cothromachd/snippetbox/pkg/models"
)

// Обработчик главной странице.

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        app.notFound(w)
        return
    }
 
    s, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, err)
        return
    }
 
    // Создаем экземпляр структуры templateData,
    // содержащий срез с заметками.
    data := &templateData{Snippets: s}
 
    files := []string{
        "./ui/html/home.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }
 
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err)
        return
    }
 
    // Передаем структуру templateData в шаблонизатор.
    // Теперь она будет доступна внутри файлов шаблона через точку.
    err = ts.Execute(w, data)
    if err != nil {
        app.serverError(w, err)
    }
}

// Обработчик для отображения содержимого заметки.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := &templateData{Snippet: s}

	files := []string {
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)

	if err != nil {
		app.serverError(w, err)
	}

}

// Обработчик для создания новой заметки.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)

		return

	}

	title := "История про улитку"
	content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return 
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
	// переопределить права web'у
	// удалить лишние ряды из таблицы snippets
}
