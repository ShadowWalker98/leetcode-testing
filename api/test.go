package main

import (
	"Testing/internal/data"
	"bufio"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type Application struct {
	cfg    config
	models data.Models
	logger *log.Logger
}

func main() {
	var cfg config

	cfg.db.dsn = "postgres://problemsuser:password@localhost/problems?sslmode=disable"

	conn, err := openDB(cfg)

	if err != nil {
		fmt.Println(err)
	}

	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := Application{
		cfg:    cfg,
		models: data.NewModels(conn),
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 3000),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
	//app.DisplayMenu()
}

func (app *Application) DisplayMenu() {
	for {
		fmt.Println("Welcome to the menu")
		fmt.Println("1. Add a problem")
		fmt.Println("2. Display problems")
		fmt.Println("3. Delete problem")
		fmt.Println("0. Exit")

		var option int64

		_, err := fmt.Scanln(&option)
		if err != nil {
			fmt.Println("Error occurred while reading input", err)
			fmt.Println()
		}

		switch option {
		case 1:
			app.AddProblemHandler()
			break
		case 2:
			app.DisplayProblemHandler()
			break
		case 3:
			app.DeleteProblemHandler()
		case 0:
			app.ExitHandler()

		default:
			fmt.Println("Please enter valid input")
		}
	}
}

func (app *Application) DeleteProblemHandler() {
	var problemNumber int

	fmt.Println("Enter problem number to delete")

	_, err := fmt.Scanln(&problemNumber)
	if err != nil {
		return
	}

	deleted := app.models.Problems.DeleteProblem(problemNumber)

	if deleted {
		fmt.Println("Deleted problem with number: ", problemNumber)
	}

}

func (app *Application) AddProblemHandler() {

	var problemNumber int
	var problemName string
	var days int
	var months int

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Adding problem")
	fmt.Println("Enter problem number")
	_, err := fmt.Scanln(&problemNumber)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Incorrect input entered, returning to menu")
		return
	}
	fmt.Println("Enter problem name")
	problemName, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		fmt.Println("Incorrect input entered, returning to menu")
		return
	}
	problemName = problemName[:len(problemName)-2]
	fmt.Println("Enter days")
	_, err = fmt.Scanln(&days)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Incorrect input entered, returning to menu")
		return
	}

	fmt.Println("Enter months")
	_, err = fmt.Scanln(&months)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Incorrect input entered, returning to menu")
		return
	}

	problem := data.Problem{
		ProblemNumber:     problemNumber,
		ProblemName:       problemName,
		LastSolvedOn:      time.Now(),
		DueDate:           time.Now().AddDate(0, months, days),
		NumberTimesSolved: 1,
	}

	existingProblem, found := app.models.Problems.SelectRowWithProblemNumber(problem.ProblemNumber)

	if found {
		fmt.Println("Problem already exists. Updating the problem ", existingProblem.ProblemNumber)
		updateErr, updated := app.models.Problems.UpdateProblem(existingProblem, problem)

		if !updated {
			fmt.Println(updateErr)
		}

		return
	} else {
		err = app.models.Problems.Insert(&problem)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Row inserted successfully")
			fmt.Println()
		}
	}
}

func (app *Application) DisplayProblemHandler() {
	fmt.Println("Displays problems solved and their information")
	app.models.Problems.ViewAllProblems()
}

func (app *Application) ExitHandler() {
	fmt.Println("Exiting program")
	os.Exit(0)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(time.Minute * 15)
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type database struct {
	dsn string
}

type config struct {
	db database
}
