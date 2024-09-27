package expense

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetExpenseByIdHandler(c echo.Context) error {
	id := c.Param("id")

	rowId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "id should be int " + err.Error()})
	}

	exp := Expense{}
	h.DB.First(&exp, rowId)

	return c.JSON(http.StatusOK, exp)
}

func (h *handler) GetExpensesHandler(c echo.Context) error {
	exps := []Expense{}
	h.DB.Find(&exps)
	return c.JSON(http.StatusOK, exps)
}
