package expense

// unit test
import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	expenseJson = bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)
)

func TestExpenseModelNotNil(t *testing.T) {
	exp := &Expense{
		Id:     1,
		Title:  "buy a new phone",
		Amount: 39000,
		Note:   "buy a new phone",
		Tags:   []string{"gadget", "shopping"},
	}

	if assert.NotNil(t, exp) {
		assert.Equal(t, 1, exp.Id)
		assert.Equal(t, "buy a new phone", exp.Title)
		assert.Equal(t, 39000.00, exp.Amount)
		assert.Equal(t, "buy a new phone", exp.Note)
		assert.Equal(t, []string{"gadget", "shopping"}, exp.Tags)
	}
}

func TestExpenseCreate(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", expenseJson)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	db, mock, err := sqlmock.New()
	mock.ExpectQuery(
		"INSERT INTO expenses \\(title, amount, note, tags\\) values \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id").WithArgs(
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg()).WillReturnRows(newsMockRows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := handler{db}
	c := e.NewContext(req, rec)

	err = h.CreateExpenseHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}