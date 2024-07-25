package data

import "database/sql"

type Models struct {
	Problems ProblemModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Problems: ProblemModel{DB: db},
	}
}
