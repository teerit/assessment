package expense

import (
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
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, exp)
}
