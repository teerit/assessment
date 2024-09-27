package expense

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) UpdateExpenseHandler(c echo.Context) error {
	rowId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	updatedExp := Expense{}
	err = c.Bind(&updatedExp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	exp := Expense{}
	h.DB.First(&exp, rowId)

	exp.Id = rowId
	exp.Title = updatedExp.Title
	exp.Amount = updatedExp.Amount
	exp.Note = updatedExp.Note
	exp.Tags = updatedExp.Tags

	h.DB.Save(&exp)

	return c.JSON(http.StatusOK, exp)
}
