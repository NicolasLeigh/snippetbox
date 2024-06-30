package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.linze.me/internal/models"
)

// Change the signature of the home handler so it is defined as a method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because httprouter matches the "/" path exactly, we can now remove the manual check of r.URL.Path != "/" from this handler.
	/*  
	// Check if the current request URL path exactly matches "/". If it doesn't, use
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the handler
	// would keep executing and also write the "Hello from SnippetBox" message.
	if r.URL.Path != "/" {
		// http.NotFound(w, r)
		app.notFound(w) // Use the notFound() helper
		return
	}
	*/

	// panic("something went wrong!")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w,err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containing the 'default' data (which for now is just the current year), and add the snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Use the new render helper

	// Pass the data to the render() helper as normal.
	app.render(w, http.StatusOK, "home.tmpl", data)

	/*
	// Initialize a slice containing the paths to the two files. It's important
	// to note that the file containing our base template must be the *first*
	// file in the slice.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and use
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file
	// paths as a variadic parameter
	ts, err := template.ParseFiles(files...)
	if err != nil {
		// Because the home handler function is now a method against application
		// it can access its fields, including the error logger. We'll write the log
		// message to this instead of the standard logger.
		// app.errLog.Print(err.Error())
		// http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err) // Use the serverError() helper
		return
	}

	// Create an instance of a templateData struct holding the slice of snippets.
	data := &templateData{
		Snippets: snippets,
	}

	// We then use the Execute() method on the template set to write the
	// template content as the response body. The last parameter to Execute()
	// represents any dynamic data that we want to pass in, which for now we'll
	// leave as nil.
	//err = ts.Execute(w, nil)

	// Use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body.
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		// Also update the code here to use the error logger from the application struct.
		// app.errLog.Print(err.Error())
		// http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err) // Use the serverError() helper
	}

	*/
}

// Change the signature of the snippetView handler so it is defined as a method
// against *application
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// When httprouter is parsing a request, the values of any named parameters will be stored in the request context. 
	// We'll talk about request context in detail later in the book, but for now it's enough to know that 
	// you can use the ParamsFromContext() function to retrieve a slice containing these parameter names and values like so:
	params := httprouter.ParamsFromContext(r.Context())

	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() function. If it can't
	// be converted to an integer, or the value is less than 1, we return a 404 page
	// not found response.

	// We can then use the ByName() method to get the value of the "id" named parameter from the slice and validate it as normal.
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFound(w) // Use the notFound() helper
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a     
	// specific record based on its ID. If no matching record is found, return a 404 Not Found response.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord){
			app.notFound(w)
		} else {
			app.serverError(w,err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	// Use the new render helper
	app.render(w, http.StatusOK, "view.tmpl", data)

	/*
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/view.tmpl",
		"./ui/html/partials/nav.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w,err)
		return
	}

	// Create an instance of a templateData struct holding the snippet data.
	data := &templateData{
		Snippet: snippet,
	}

	// Notice how we are passing in the snippet data (a models.Snippet struct) as the final parameter

	// Pass in the templateData struct when executing the template.
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w,err)
	}
	// Use the fmt.Fprintf() function to interpolate the id value with our response
	// and write it to the http.ResponseWriter.
	// fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)

	// Write the snippet data as a plain-text HTTP response body.
	// The plus flag (%+v) adds field names
	fmt.Fprintf(w, "%+v", snippet)

	*/
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "create.tmpl", data)
}

// Change the signature of the snippetCreate handler so it is defined as a method
// against *application.

// Rename this handler to snippetCreatePost.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Checking if the request method is a POST is now superfluous and can be removed, because this is done automatically by httprouter.
	/* 
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper
		return
	}
	*/

	// First we call r.ParseForm() which adds any data in POST request bodies to the r.PostForm map. This also works in the same way for PUT and PATCH requests. 
	// If there are any errors, we use our app.ClientError() helper to send a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the title and content from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// The r.PostForm.Get() method always returns the form data as a *string*.
  // However, we're expecting our expires value to be a number, and want to represent it in our Go code as an integer. 
	// So we need to manually covert the form data to an integer using strconv.Atoi(), and we send a 400 Bad Request response if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return 
	}

	/* 
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7
	*/

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	// Update the redirect path to use the new clean URL format.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}

func (app *application) render (w http.ResponseWriter, status int, page string, data *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name (like 'home.tmpl'). 
	// If no entry exists in the cache with the provided name, then create a new error and call the serverError() helper method that we made earlier and return.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the http.ResponseWriter. 
	// If there's an error, call our serverError() helper and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer without any errors, 
	// we are safe to go ahead and write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter. 
	// Note: this is another time where we pass our http.ResponseWriter to a function that takes an io.Writer.
	buf.WriteTo(w)

 /* 
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
*/
}