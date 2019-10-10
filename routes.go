package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

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
		"Healthy",
		"GET",
		"/api/liquidacion/healthy",
		Healthy,
	},
	Route{
		"LiquidacionList",
		"GET",
		"/api/liquidacion/liquidaciones",
		LiquidacionList,
	},
	Route{
		"LiquidacionShow",
		"GET",
		"/api/liquidacion/liquidaciones/{id}",
		LiquidacionShow,
	},
	Route{
		"LiquidacionAdd",
		"POST",
		"/api/liquidacion/liquidaciones",
		LiquidacionAdd,
	},
	Route{
		"LiquidacionUpdate",
		"PUT",
		"/api/liquidacion/liquidaciones/{id}",
		LiquidacionUpdate,
	},
	Route{
		"LiquidacionRemove",
		"DELETE",
		"/api/liquidacion/liquidaciones/{id}",
		LiquidacionRemove,
	},
	Route{
		"LiquidacionContabilizar",
		"POST",
		"/api/liquidacion/contabilizar",
		LiquidacionContabilizar,
	},
	/*Route{
		"LiquidacionDesContabilizar",
		"POST",
		"/api/liquidacion/descontabilizar",
		LiquidacionDesContabilizar,
	},*/
	Route{
		"LiquidacionesRemoveMasivo",
		"DELETE",
		"/api/liquidacion/liquidaciones",
		LiquidacionesRemoveMasivo,
	},
	Route{
		"LiquidacionDuplicarMasivo",
		"POST",
		"/api/liquidacion/duplicar",
		LiquidacionDuplicarMasivo,
	},
}
