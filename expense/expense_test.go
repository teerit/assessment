package expense

// unit test
import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	expenseJson = `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`
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
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(expenseJson))
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
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestExpenseGetById(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// mock result
	mockTags := []string{"food", "beverage"}
	mockRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "strawberry smoothie", "79", "night market promotion discount 10 bath", pq.Array(&mockTags))
	// expected
	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}\n"

	db, mock, err := sqlmock.New()
	mock.ExpectQuery(
		"SELECT (.+) FROM expenses WHERE id=\\$1").WithArgs(sqlmock.AnyArg()).WillReturnRows(mockRow)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	err = h.GetExpenseByIdHandler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
	}
}

func TestExpenseUpdateById(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(expenseJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	id := "1"

	mockTags := []string{"food", "beverage"}

	// expected
	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}\n"

	db, mock, err := sqlmock.New()
	stmt := mock.ExpectPrepare("UPDATE expenses SET title=\\$2, amount=\\$3, note=\\$4, tags=\\$5 WHERE id=\\$1")
	stmt.ExpectExec().
		WithArgs(
			1,
			"strawberry smoothie",
			79.00,
			"night market promotion discount 10 bath",
			pq.Array(mockTags)).WillReturnResult(sqlmock.NewResult(1, 1))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues(id)
	err = h.UpdateExpenseHandler(c)

	// assertion
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
	}
}

func TestExpenseGetAll(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// mock result
	mockTags := []string{"food", "beverage"}
	mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow("1", "strawberry smoothie", "79", "night market promotion discount 10 bath", pq.Array(&mockTags))

	// expected return
	expected := "[{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}]"
	db, mock, err := sqlmock.New()
	mock.ExpectQuery("SELECT (.+) FROM expenses").WillReturnRows(mockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := handler{db}
	c := e.NewContext(req, rec)

	err = h.GetExpensesHandler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
