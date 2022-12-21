package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/teerit/assessment/db"
	"github.com/teerit/assessment/expense"
	"github.com/teerit/assessment/middleware"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		fmt.Printf("Error initial db connection %s", err)
	}

	h := expense.ExpenseHandler(db)
	e := echo.New()

	e.Use(middleware.DateFormatAuthMiddleware)
	e.Use(middleware.RequestLogger)

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
