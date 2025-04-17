package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Nikita-Astafyev/book-service/models"
	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	DB *sql.DB
}

func NewBookHandler(db *sql.DB) *BookHandler {
	return &BookHandler{DB: db}
}

func (h *BookHandler) CreateBook(c echo.Context) error {
	var book models.Book
	if err := c.Bind(&book); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if book.ID == "" || book.Title == "" || book.Author == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "All fields are required")
	}

	_, err := h.DB.Exec(
		"INSERT INTO books(id, title, author) VALUES($1, $2, $3)",
		book.ID, book.Title, book.Author,
	)

	if err != nil {
		log.Printf("Database error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) GetBook(c echo.Context) error {
	id := c.Param("id")

	var book models.Book
	row := h.DB.QueryRow("SELECT id, title, author FROM books WHERE id = $1", id)
	err := row.Scan(&book.ID, &book.Title, &book.Author)

	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, "Book not found")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, book)
}

func (h *BookHandler) UpdateBook(c echo.Context) error {
	id := c.Param("id")

	var book models.Book
	if err := c.Bind(&book); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if book.Title == "" || book.Author == "" {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Title and authors are required fields",
		)
	}

	result, err := h.DB.Exec(
		"UPDATE books SET title = $1, author = $2 WHERE id = $3",
		book.Title, book.Author, book.ID,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			"Database update error", err.Error(),
		)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return echo.NewHTTPError(
			http.StatusNotFound,
			"Book with id"+id+"not found",
		)
	}

	return c.JSON(http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(c echo.Context) error {
	id := c.Param("id")

	result, err := h.DB.Exec("DELETE from books WHERE id = $1", id)

	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Database delete error", err.Error(),
		)
	}

	rowAffected, _ := result.RowsAffected()
	if rowAffected == 0 {
		return echo.NewHTTPError(
			http.StatusNotFound,
			"Book with id", id, "not found",
		)
	}

	return c.NoContent(http.StatusNoContent)
}
