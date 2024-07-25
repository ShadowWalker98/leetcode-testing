package data

import (
	"Testing/internal/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var ProblemNumberKey = "problem number"
var DueDateKey = "due date"
var NumberTimesSolvedKey = "times solved"

type Problem struct {
	ProblemNumber     int
	ProblemName       string
	LastSolvedOn      time.Time
	DueDate           time.Time
	NumberTimesSolved int
}

func (p Problem) ToString() string {
	str := strconv.Itoa(p.ProblemNumber) + " " + p.ProblemName
	return str
}

type ProblemModel struct {
	DB *sql.DB
}

func (p ProblemModel) Insert(problem *Problem) error {
	query := `INSERT INTO problems (problem_number, problem_name, last_solved_on, due_date, number_times_solved)
			  VALUES ($1, $2, $3, $4, $5)`
	args := []interface{}{problem.ProblemNumber, problem.ProblemName, problem.LastSolvedOn, problem.DueDate, problem.NumberTimesSolved}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, args...)
	if err != nil {
		fmt.Println(err)
		return errors.New("error while insertion")
	}
	return nil
}

func (p ProblemModel) SelectRowWithProblemNumber(problemNumber int) (Problem, bool) {
	query := `SELECT * FROM problems WHERE problem_number = $1`
	args := []interface{}{problemNumber}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var problem Problem

	err := p.DB.QueryRowContext(ctx, query, args...).Scan(
		&problem.ProblemNumber,
		&problem.ProblemName,
		&problem.LastSolvedOn,
		&problem.DueDate,
		&problem.NumberTimesSolved)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("No rows with corresponding problem number found")
		} else {
			fmt.Println(err)
		}

		return problem, false
	}

	return problem, true
}

func (p ProblemModel) UpdateProblem(oldProblem Problem, newProblem Problem) (error, bool) {
	query := `UPDATE problems SET last_solved_on = $1, due_date = $2, number_times_solved = $3;`
	args := []interface{}{newProblem.LastSolvedOn, newProblem.DueDate, oldProblem.NumberTimesSolved + 1}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := p.DB.ExecContext(ctx, query, args...)

	if err != nil {
		fmt.Println(err)
		return errors.New("error while updating row"), false
	}
	fmt.Printf("Rows updated: %d\n", rows)
	return nil, true
}

func (p ProblemModel) ViewAllProblemsResponseWriter(w http.ResponseWriter) {
	query := `SELECT * FROM problems`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	rows, err := p.DB.QueryContext(ctx, query)

	for rows.Next() {
		var problem Problem
		err := rows.Scan(&problem.ProblemNumber,
			&problem.ProblemName,
			&problem.LastSolvedOn,
			&problem.DueDate,
			&problem.NumberTimesSolved)
		if err != nil {
			fmt.Println(err)

		}
		_, err = fmt.Fprint(w, problem.ToString()+"\n")
		if err != nil {
			fmt.Println(err)
		}
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("No problems solved.")
		}
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
}

func (p ProblemModel) ViewAllProblems() {

	query := `SELECT * FROM problems`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	rows, err := p.DB.QueryContext(ctx, query)

	for rows.Next() {
		var problem Problem
		err := rows.Scan(&problem.ProblemNumber,
			&problem.ProblemName,
			&problem.LastSolvedOn,
			&problem.DueDate,
			&problem.NumberTimesSolved)
		if err != nil {
			fmt.Println(err)

		}
		fmt.Println(problem.ToString())
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("No problems solved.")
		}
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
}

func (p ProblemModel) DeleteProblem(problemNumber int) bool {
	query := `DELETE FROM problems WHERE problem_number = $1`
	args := []interface{}{problemNumber}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, args...)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
func ValidateProblemData(v *validator.Validator, problem *Problem) {
	v.Check(problem.ProblemNumber >= 1, ProblemNumberKey, "problem number must be greater than 0")
	v.Check(problem.DueDate.After(time.Now()), DueDateKey, "due date must be later than the current time")
	v.Check(problem.NumberTimesSolved >= 0, NumberTimesSolvedKey, "number of times must be >=0")
}
