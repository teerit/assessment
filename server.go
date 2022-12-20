package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/teerit/assessment/expense"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	h := expense.ExpenseHandler(db)
	createTable := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`

	_, err = h.DB.Exec(createTable)
	if err != nil {
		log.Fatal("can't create table", err)
	}

	e := echo.New()
	e.POST("/expenses", h.CreateExpenseHandler)
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
