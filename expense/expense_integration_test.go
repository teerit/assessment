//go:build integration
// +build integration

package expense

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/teerit/assessment/util"
)

func TestITExpenses(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/assessment-db?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h := ExpenseHandler(db)

		e.POST("/expenses", h.CreateExpenseHandler)
		e.GET("/expenses/:id", h.GetExpenseByIdHandler)
		e.GET("/expenses", h.GetExpensesHandler)
		e.PUT("/expenses/:id", h.UpdateExpenseHandler)

		e.Start(fmt.Sprintf(":%d", util.ServerPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", util.ServerPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	t.Run("TestGetExpenses", func(t *testing.T) {
		seedExpense(t)
		var exps []Expense

		res := util.Request(http.MethodGet, util.Uri("expenses"), nil)
		err := res.Decode(&exps)

		assert.Nil(t, err)
		assert.EqualValues(t, http.StatusOK, res.StatusCode)
		assert.Greater(t, len(exps), 0)
	})

	t.Run("TestGetExpenseById", func(t *testing.T) {
		c := seedExpense(t)

		var lastExp Expense
		res := util.Request(http.MethodGet, util.Uri("expenses", strconv.Itoa(c.Id)), nil)
		err := res.Decode(&lastExp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, c.Id, lastExp.Id)
		assert.NotEmpty(t, lastExp.Title)
		assert.NotEmpty(t, lastExp.Amount)
		assert.NotEmpty(t, lastExp.Note)
	})

	t.Run("TestUpdateExpense", func(t *testing.T) {
		id := seedExpense(t).Id
		c := Expense{
			Id:     id,
			Title:  "strawberry smoothie",
			Amount: 79,
			Note:   "night market promotion discount 10 bath",
			Tags:   []string{"food", "beverage"},
		}
		payload, _ := json.Marshal(c)
		res := util.Request(http.MethodPut, util.Uri("expenses", strconv.Itoa(id)), bytes.NewBuffer(payload))
		var info Expense
		err := res.Decode(&info)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, c.Title, info.Title)
		assert.Equal(t, c.Amount, info.Amount)
		assert.Equal(t, c.Note, info.Note)
	})

	t.Run("TestCreateExpense", func(t *testing.T) {
		body := bytes.NewBufferString(`{
			"title": "strawberry smoothie",
			"amount": 79,
			"note": "night market promotion discount 10 bath", 
			"tags": ["food", "beverage"]
		}`)
		var exp Expense

		res := util.Request(http.MethodPost, util.Uri("expenses"), body)
		err := res.Decode(&exp)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.NotEqual(t, 0, exp.Id)
		assert.Equal(t, "strawberry smoothie", exp.Title)
		assert.Equal(t, 79.0, exp.Amount)
		assert.Equal(t, "night market promotion discount 10 bath", exp.Note)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func seedExpense(t *testing.T) Expense {
	var c Expense
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)
	err := util.Request(http.MethodPost, util.Uri("expenses"), body).Decode(&c)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}
	return c
}
