package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// app.errLog.Print(trace)
	app.errLog.Output(2, trace) // calldepth is 2 // ERROR   2024/06/25 13:25:25 handlers.go:46: open ./ui/html/pages/home.tmpl: no such file or directory
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Create an newTemplateData() helper, which returns a pointer to a templateData struct initialized with the current year. 
//Note that we're not using the *http.Request parameter here at the moment, but we will do later in the book.
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{CurrentYear: time.Now().Year(),
	// Add the flash message to the template data, if one exists.
	Flash: app.sessionManager.PopString(r.Context(), "flash"),
	}
}