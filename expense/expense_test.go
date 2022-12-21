//go:build unit
// +build unit

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

func testWrapper(jsonString string) (*http.Request, *httptest.ResponseRecorder, *echo.Echo) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(jsonString))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return req, rec, e
}

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

func TestExpenseHandler(t *testing.T) {
	db, _, err := sqlmock.New()
	con := ExpenseHandler(db)
	if assert.NoError(t, err) {
		assert.NotNil(t, con)
	}
}

func TestExpenseCreate(t *testing.T) {
	tests := []struct {
		name         string
		expectedCode int
		mockRows     *sqlmock.Rows
	}{
		{
			name:         "TestExpenseCreateSuccess",
			expectedCode: http.StatusCreated,
			mockRows:     sqlmock.NewRows([]string{"id"}).AddRow("1"),
		},
		{
			name:         "TestExpenseCreateInternalServerError",
			expectedCode: http.StatusInternalServerError,
			mockRows:     sqlmock.NewRows([]string{"id"}).AddRow("xxx"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, rec, e := testWrapper(expenseJson)
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			mock.ExpectQuery(
				"INSERT INTO expenses \\(title, amount, note, tags\\) values \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id").
				WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg()).WillReturnRows(test.mockRows)

			h := handler{db}
			c := e.NewContext(req, rec)

			err = h.CreateExpenseHandler(c)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expectedCode, rec.Code)
			}
		})
	}
}

func TestExpenseGetById(t *testing.T) {
	tests := []struct {
		name         string
		paramValue   string
		expectedCode int
		expectedBody string
		mockRows     *sqlmock.Rows
	}{
		{
			name:         "TestExpenseGetSuccess",
			paramValue:   "1",
			expectedCode: http.StatusOK,
			expectedBody: "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}\n",
			mockRows: sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
				AddRow("1", "strawberry smoothie", "79", "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"})),
		},
		{
			name:         "TestExpenseGetNotFound",
			paramValue:   "1",
			expectedCode: http.StatusNotFound,
			expectedBody: "{\"message\":\"expense not found with given id\"}\n",
			mockRows:     sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, rec, e := testWrapper("")
			db, mock, err := sqlmock.New()

			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			mock.ExpectQuery(
				"SELECT (.+) FROM expenses WHERE id=\\$1").WithArgs(sqlmock.AnyArg()).
				WillReturnRows(test.mockRows)

			h := handler{db}
			c := e.NewContext(req, rec)
			c.SetPath("/expenses/:id")
			c.SetParamNames("id")
			c.SetParamValues(test.paramValue)
			err = h.GetExpenseByIdHandler(c)

			if assert.NoError(t, err) {
				assert.Equal(t, test.expectedCode, rec.Code)
				assert.Equal(t, test.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestExpenseUpdateById(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		pathParam      string
		tags           []string
		expected       string
		expectedStatus int
	}{
		{
			name:           "TestExpenseUpdateSuccess",
			requestBody:    expenseJson,
			pathParam:      "1",
			tags:           []string{"food", "beverage"},
			expected:       "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}\n",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "TestExpenseUpdateBadRequest",
			requestBody:    expenseJson,
			pathParam:      "",
			tags:           []string{"food", "beverage"},
			expected:       "{\"message\":\"strconv.Atoi: parsing \\\"\\\": invalid syntax\"}\n",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, rec, e := testWrapper(test.requestBody)

			db, mock, err := sqlmock.New()
			stmt := mock.ExpectPrepare("UPDATE expenses SET title=\\$2, amount=\\$3, note=\\$4, tags=\\$5 WHERE id=\\$1")
			stmt.ExpectExec().
				WithArgs(
					1,
					"strawberry smoothie",
					79.00,
					"night market promotion discount 10 bath",
					pq.Array(test.tags)).WillReturnResult(sqlmock.NewResult(1, 1))

			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			h := handler{db}
			c := e.NewContext(req, rec)
			c.SetPath("/expenses/:id")
			c.SetParamNames("id")
			c.SetParamValues(test.pathParam)
			err = h.UpdateExpenseHandler(c)

			// assertion
			if assert.NoError(t, err) {
				assert.Equal(t, test.expectedStatus, rec.Code)
				assert.Equal(t, test.expected, rec.Body.String())
			}
		})
	}
}

func TestExpenseGetAll(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		tags           []string
		expected       string
		expectedStatus int
	}{
		{
			name:           "TestExpenseGetAllSuccess",
			requestBody:    "",
			tags:           []string{"food", "beverage"},
			expected:       "[{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}]",
			expectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, rec, e := testWrapper(test.requestBody)

			mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
				AddRow("1", "strawberry smoothie", "79", "night market promotion discount 10 bath", pq.Array(&test.tags))

			db, mock, err := sqlmock.New()
			mock.ExpectQuery("SELECT (.+) FROM expenses").WillReturnRows(mockRows)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			h := handler{db}
			c := e.NewContext(req, rec)

			err = h.GetExpensesHandler(c)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expectedStatus, rec.Code)
				assert.Equal(t, test.expected, strings.TrimSpace(rec.Body.String()))
			}
		})
	}
}
