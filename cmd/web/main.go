package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"     
	// and some short help text explaining what the flag controls. The value of the     
	// flag will be stored in the addr variable at runtime
	addr := flag.String("addr",":4000","Http network address")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
  // This reads in the command-line flag value and assigns it to the addr     
	// variable. You need to call this *before* you use the addr variable     
	// otherwise it will always contain the default value of ":4000". If any errors are     
	// encountered during parsing the application will be terminated
	flag.Parse()

	// Use the http.NewServeMux() function to initialize a new servemux, then     
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
  // Note that the path given to the http.Dir function is relative to the project     
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/",http.StripPrefix("/static/",fileServer))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Use the http.ListenAndServe() function to start a new web server. We pass in     
	// two parameters: the TCP network address to listen on (in this case ":4000")     
	// and the servemux we just created. If http.ListenAndServe() returns an error     
	// we use the log.Fatal() function to log the error message and exit. Note     
	// that any error returned by http.ListenAndServe() is always non-nil.

	// The value returned from the flag.String() function is a pointer to the flag     
	// value, not the value itself. So we need to dereference the pointer (i.e.
  // prefix it with the * symbol) before using it. Note that we're using the     
	// log.Printf() function to interpolate the address with the log message.
	log.Printf("Starting server on %s",*addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
