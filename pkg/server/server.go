package server

import (
	"log"
	"net/http"

	"github.com/dafiti-group/aws-s3-sync-api/pkg/handler"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
)

// Server has router and db instances
type Server struct {
	Log     logr.Logger
	Handler handler.Handler
	*mux.Router
}

// Initialize initializes the app with predefined configuration
func (a *Server) Initialize() {
	a.Router = mux.NewRouter()
	a.Handler = handler.Handler{
		Log: a.Log.WithName("server"),
	}
	a.setRouters()
}

// setRouters sets the all required routers
func (a *Server) setRouters() {
	// Routing for handling the projects
	a.Get("/{path:.+}", a.handleRequest(a.Handler.GetAllFiles))
	a.Post("/", a.handleRequest(a.Handler.Sync))
}

// Get wraps the router for GET method
func (a *Server) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *Server) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Run the app on it's router
func (a *Server) Run(host string) {
	a.Log.WithName("server").Info("Server Started")
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)

func (a *Server) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}
