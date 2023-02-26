package expense

import (
	"fmt"
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
	if result := h.DB.First(&exp, rowId); result.Error != nil {
		fmt.Println(result.Error)
	}

	exp.Title = updatedExp.Title
	exp.Amount = updatedExp.Amount
	exp.Note = updatedExp.Note
	exp.Tags = updatedExp.Tags

	h.DB.Save(&exp)

	return c.JSON(http.StatusOK, exp)
}
