package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.

// Update the signature for the routes() method so that it returns a http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler{
	/* 
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.
  // Because secureHeaders is just a function, and the function returns a http.Handler we don't need to do anything else.
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	*/

	// Initialize the router
	router := httprouter.New()

	// Create a handler function which wraps our notFound() helper, and then assign it as the custom handler for 404 Not Found responses. 
	// You can also set a custom handler for 405 Method Not Allowed responses by setting router.MethodNotAllowed in the same way too.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Update the pattern for the route for the static files.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static/", fileServer))

	// Create a new middleware chain containing the middleware specific to our dynamic application routes. For now, this chain will only contain the LoadAndSave session middleware but we'll add more to it later.
	// Unprotected application routes using the "dynamic" middleware chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	/* 
	// And then create the routes using the appropriate methods, patterns and handlers.
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)
 */

	// Update these routes to use the new dynamic middleware chain followed by the appropriate handler function. 
	// Note that because the alice ThenFunc() method returns a http.Handler (rather than a http.HandlerFunc) we also need to switch to registering the route using the router.Handler() method.
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected (authenticated-only) application routes, using a new "protected" middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// Create the middleware chain as normal.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Wrap the router with the middleware and return it as normal.
	return standard.Then(router)
}
