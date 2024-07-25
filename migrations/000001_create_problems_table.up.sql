CREATE TABLE IF NOT EXISTS problems
(
    problem_number      int PRIMARY KEY,
    problem_name        text,
    last_solved_on      date NOT NULL,
    due_date            date NOT NULL,
    number_times_solved int  NOT NULL DEFAULT 0
);

