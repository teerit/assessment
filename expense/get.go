package expense

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) GetExpenseByIdHandler(c echo.Context) error {
	id := c.Param("id")

	rowId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "id should be int " + err.Error()})
	}

	row := h.DB.QueryRow("SELECT id, title, amount, note, tags FROM expenses WHERE id=$1", rowId)

	exp := Expense{}
	err = row.Scan(&exp.Id, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, Err{Message: "expense not found with given id"})
		}
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, exp)
}

func (h *handler) GetExpensesHandler(c echo.Context) error {
	exps := []Expense{}

	rows, err := h.DB.Query("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	for rows.Next() {
		exp := Expense{}
		err := rows.Scan(&exp.Id, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})

		}
		exps = append(exps, exp)
	}

	return c.JSON(http.StatusOK, exps)
}
