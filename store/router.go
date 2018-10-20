package store

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
)


var controller = &Controller{Repository: Repository{}}

// Route defines a route
type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

// Route defines the list of routes of API
type Routes []Route

var routes = Routes {
	Route {
		"Authentication",
		"POST",
		"/get-token",
		controller.GetToken,
	},
	Route {
		"Index",
		"GET",
		"/",
		controller.Index,
	},
	Route {
		"AddProduct",
		"POST",
		"/AddProduct",
		AuthenticationMiddleware(controller.AddProduct),
	},
	// More routes..
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		log.Println(route.Name)
		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}