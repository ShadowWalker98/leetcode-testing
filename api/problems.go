package main

import (
	"fmt"
	"net/http"
)

func (app *Application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Println("healthcheck working")
	_, err := fmt.Fprintf(w, "healthcheck info")
	if err != nil {
		return
	}
}

func (app *Application) DisplayProblemsHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Println("display problems")
	app.models.Problems.ViewAllProblemsResponseWriter(w)
	//for rowsNext() {
	//	var problem data.Problem
	//	err := rows.Scan(&problem.ProblemNumber,
	//		&problem.ProblemName,
	//		&problem.LastSolvedOn,
	//		&problem.DueDate,
	//		&problem.NumberTimesSolved)
	//
	//	if err != nil {
	//		panic(err)
	//	}
	//	app.logger.Println(problem.ToString())
	//	_, err = w.Write([]byte(problem.ToString()))
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//}
}
