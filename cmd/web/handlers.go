package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	// "strings"
	// "unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"snippetbox.linze.me/internal/models"
	"snippetbox.linze.me/internal/validator"
)

// Define a snippetCreateForm struct to represent the form data and validation errors for the form fields. Note that all the struct fields are deliberately exported (i.e. start with a capital letter).
// This is because struct fields must be exported in order to be read by the html/template package when rendering the template.

// Remove the explicit FieldErrors struct field and instead embed the Validator type.
// Embedding this means that our snippetCreateForm "inherits" all the fields and methods of our Validator type (including the FieldErrors field).
type snippetCreateForm struct {
	Title string
	Content string
	Expires int
	// FieldErrors map[string]string
	validator.Validator
}

type userSignupForm struct {
	Name string `form:"name"`
	Email string `form:"email"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
}

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

	// Use the PopString() method to retrieve the value for the "flash" key.
  // PopString() also deletes the key and value from the session data, so it acts like a one-time fetch. 
	// If there is no matching key in the session data this will return the empty string.
	// No need to add this line of code, because we are using the newTemplateData() helper to do this for us.
	// flash := app.sessionManager.PopString(r.Context(), "flash")

	data := app.newTemplateData(r)
	data.Snippet = snippet
	// Same as before, we pass the flash message to the template data.
	// data.Flash = flash

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

	// Initialize a new createSnippetForm instance and pass it to the template.
  // Notice how this is also a great opportunity to set any default or 'initial' values for the form --- here we set the initial value for the snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		Expires: 365,
	}
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
/* 
	// Use the r.PostForm.Get() method to retrieve the title and content from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
*/

	// The r.PostForm.Get() method always returns the form data as a *string*.
  // However, we're expecting our expires value to be a number, and want to represent it in our Go code as an integer. 
	// So we need to manually covert the form data to an integer using strconv.Atoi(), and we send a 400 Bad Request response if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return 
	}

	// Create an instance of the snippetCreateForm struct containing the values from the form and an empty map for any validation errors.
	form := snippetCreateForm{
		Title: r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
		// FieldErrors: map[string]string{},
	}

	// Because the Validator type is embedded by the snippetCreateForm struct, we can call CheckField() directly on it to execute our validation checks.
  // CheckField() will add the provided key and error message to the  FieldErrors map if the check does not evaluate to true. 
	// For example, in the first line here we "check that the form.Title field is not blank". 
	// In the second, we "check that the form.Title field has a maximum character length of 100" and so on.
	form.CheckField(validator.NotBlank(form.Title), "title", "must not be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "must not be more than 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "must not be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "must be a valid expiry period")

	// Use the Valid() method to see if any of the checks failed. If they did, then re-render the template passing in the form in the same way as before.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}
	
	/* 
	// Check that the title value is not blank and is not more than 100 characters long. 
	// If it fails either of those checks, add a message to the errors map using the field name as the key.

	// Update the validation checks so that they operate on the snippetCreateForm instance.
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	// Check that the content value isn't blank.
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	// Check the expires value matches one of the permitted values (1, 7 or 365)
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	// If there are any validation errors re-display the create.tmpl template, passing in the snippetCreateForm instance as dynamic data in the Form field. 
	// Note that we use the HTTP status code 422 Unprocessable Entity when sending the response to indicate that there was a validation error.
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}
	*/
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the Put() method to add a string value ("Snippet successfully created!") and the corresponding key ("flash") to the session data.
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

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

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// err := r.ParseForm()
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }
	//  form := userSignupForm{
	// 	Name: r.PostForm.Get("name"),
	// 	Email: r.PostForm.Get("email"),
	// 	Password: r.PostForm.Get("password"),
	//  }

	// Using third-party package
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	 
	 form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	 form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	 form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	 form.CheckField(validator.NotBlank(form.Password), "password", "This filed cannot be blank")
	 form.CheckField(validator.MinChars(form.Password, 6), "password", "This filed must be at least 6 characters long")

	 if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	 }

	 // Try to create a new user record in the database. If the email already exists then add an error message to the form and re-display it.
	 err = app.users.Insert(form.Name, form.Email, form.Password)
	 if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	 }

	 // Otherwise add a confirmation flash message to the session confirming that their signup worked.
	 app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	 // And redirect the user to the login page.
	 http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}