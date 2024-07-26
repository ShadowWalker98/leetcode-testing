package main

import (
	"Testing/internal/data"
	"Testing/internal/validator"
	"net/http"
	"strings"
	"time"
)

func (app *Application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "available",
		"version": app.version,
	}

	env := envelope{
		"system_info": data,
	}

	err := app.writeJson(w, env, http.StatusOK, nil)

	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}

func (app *Application) DisplayProblemsHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Println("display problems")
	problems := app.models.Problems.ViewAllProblemsResponseWriter()
	app.writeJson(w, envelope{"problem": problems}, http.StatusOK, nil)

}

func (app *Application) DisplayProblemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.ReadIDParam(r)

	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Please provide problem number for searching problems", http.StatusBadRequest)
	}

	problem, found := app.models.Problems.ViewProblemWithNumber(int(id))

	if found {
		err := app.writeJson(w, envelope{"problem": problem}, http.StatusOK, nil)
		if err != nil {
			app.logger.Println(err)
			http.Error(w, "Problem with id "+string(id)+" doesn't exist", http.StatusNotFound)
		}
	}
}

func (app *Application) AddProblemHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProblemNumber int    `json:"problem_number"`
		ProblemName   string `json:"problem_name"`
		DueDays       int    `json:"due_days"`
		DueMonths     int    `json:"due_months"`
		DueYears      int    `json:"due_years"`
	}

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	problem := data.Problem{
		ProblemNumber:     input.ProblemNumber,
		ProblemName:       input.ProblemName,
		LastSolvedOn:      time.Now(),
		DueDate:           time.Now().AddDate(input.DueYears, input.DueMonths, input.DueDays),
		NumberTimesSolved: 1,
	}

	v := validator.New()

	if data.ValidateProblemData(v, &problem); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	existingProblem, found := app.models.Problems.SelectRowWithProblemNumber(problem.ProblemNumber)

	if found {
		problem.NumberTimesSolved = existingProblem.NumberTimesSolved + 1
		err = app.models.Problems.UpdateProblem(problem)
		if err != nil {
			app.logger.Println(err)
		}
		app.writeJson(w, envelope{"problem_added": problem}, http.StatusOK, nil)
		return
	}

	problem.ProblemName = strings.Trim(problem.ProblemName, " \n")
	// TODO: capitalise the first letter of the first word in the name

	if data.ValidateProblemData(v, &problem); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Problems.Insert(&problem)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJson(w, envelope{"problem_added": problem}, http.StatusOK, nil)
}

func (app *Application) DueProblemsHandler(w http.ResponseWriter, r *http.Request) {

	date := time.Date(time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		0,
		0,
		0,
		0,
		time.Now().Location())

	problemList := app.models.Problems.FetchProblemsDueOnOrAfter(date)

	app.writeJson(w, envelope{"problem": problemList}, http.StatusOK, nil)
}
