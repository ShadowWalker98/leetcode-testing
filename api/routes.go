package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *Application) routes() *httprouter.Router {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/addproblem", app.AddProblemHandler)
	router.HandlerFunc(http.MethodGet, "/problems", app.DisplayProblemsHandler)
	router.HandlerFunc(http.MethodGet, "/problems/:id", app.DisplayProblemHandler)
	router.HandlerFunc(http.MethodGet, "/dueproblems", app.DueProblemsHandler)

	return router
}
