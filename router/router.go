package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nekonotes/controller"
	"strings"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", controller.Homepage).Methods("GET")
	r.HandleFunc("/register", controller.Register).Methods("POST")
	r.HandleFunc("/login", controller.Login).Methods("POST")
	r.HandleFunc("/createnotes", controller.CreateNotes).Methods("POST")
	r.HandleFunc("/home",controller.GetNotesTitle).Methods(("GET"))
	r.HandleFunc("/notes",controller.GetNote).Methods(("POST"))
	r.HandleFunc("/deleteuser",controller.DeleteUser).Methods(("DELETE"))
	r.HandleFunc("/deletenote",controller.DeleteNote).Methods(("DELETE"))
	//print all the routes
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	
	return r
}
