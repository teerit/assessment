package expense

import (
	"log"
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
	if result := h.DB.First(&exp, rowId); result.Error != nil {
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found with given id"})
	}

	return c.JSON(http.StatusOK, exp)
}

func (h *handler) GetExpensesHandler(c echo.Context) error {
	exps := []Expense{}
	if result := h.DB.Find(&exps); result.Error != nil {
		log.Println(result.Error)
	}
	return c.JSON(http.StatusOK, exps)
}
