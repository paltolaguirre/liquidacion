package main

import "github.com/gorilla/mux"
import "net/http"

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

type Routes []Route

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandleFunc)

	}

	return router
}

var routes = Routes{
	Route{
		"LiquidacionList",
		"GET",
		"/liquidaciones",
		LiquidacionList,
	},
	Route{
		"LiquidacionShow",
		"GET",
		"/liquidaciones/{id}",
		LiquidacionShow,
	},
	Route{
		"LiquidacionAdd",
		"POST",
		"/liquidaciones",
		LiquidacionAdd,
	},
	Route{
		"LiquidacionUpdate",
		"PUT",
		"/liquidaciones/{id}",
		LiquidacionUpdate,
	},
	Route{
		"LiquidacionRemove",
		"DELETE",
		"/liquidaciones/{id}",
		LiquidacionRemove,
	},
}
