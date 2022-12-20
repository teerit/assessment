package expense

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) UpdateExpenseHandler(c echo.Context) error {
	rowId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	exp := Expense{}
	// find by id
	err = c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err := h.DB.Prepare(`UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1`)
	if err != nil {
		fmt.Println("ERR::", err.Error())
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if _, err := stmt.Exec(rowId, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags)); err != nil {
		fmt.Println("ERR::", err.Error())
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	exp.Id = rowId
	return c.JSON(http.StatusOK, exp)
}
