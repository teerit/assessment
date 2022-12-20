package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/teerit/assessment/expense"
)

func main() {
	db, err := initDB()
	if err != nil {
		fmt.Printf("Error initial db connection %s", err)
	}

	h := expense.ExpenseHandler(db)
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		BeforeNextFunc: func(c echo.Context) {
			c.Set("customValueFromContext", time.Now().String())
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			value, _ := c.Get("customValueFromContext").(string)
			fmt.Printf("REQUEST: uri: %v, status: %v, datetime: %v\n", v.URI, v.Status, value)
			return nil
		},
	}))

	e.POST("/expenses", h.CreateExpenseHandler)
	e.GET("/expenses/:id", h.GetExpenseByIdHandler)
	e.GET("/expenses", h.GetExpensesHandler)
	e.PUT("/expenses/:id", h.UpdateExpenseHandler)

	// Start server
	go func() {
		fmt.Println(e.Start(":" + os.Getenv("PORT")))
	}()

	// Gracefully Shutdown
	// Make channel listen for signals from OS
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	<-gracefulStop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server %s", err)
	} else {
		fmt.Println("Server gracefully stopped")
	}

	if err := h.DB.Close(); err != nil {
		fmt.Printf("Error closing db connection %s", err)
	} else {
		fmt.Println("DB connection gracefully closed")
	}
}

func initDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		return nil, fmt.Errorf("can't create table: %w", err)
	}

	return db, nil
}
