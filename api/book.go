package api

import (
  "encoding/json"
  "io/ioutil"
  "net/http"
)

// Book struct include title, author
// isbn and optional description
type Book struct {
  Title       string `json:"title"` // marshal with lower case
  Author      string `json:"author"`
  ISBN        string `json:"isbn"`
  Description string `json:"description,omitempty"`
}

// Books is a map with key equal to isbn
var books = map[string]Book{
  "0345391802": {Title: "Guide to Galaxy", Author: "Douglas Adams", ISBN: "0345391802"},
  "0123456789": {Title: "Cloud Native Go", Author: "M.-Leander Reimer", ISBN: "0123456789"},
}

// ToJSON used to marshalling of Book type
func (b Book) ToJSON() []byte {
  ToJSON, err := json.Marshal(b)
  if err != nil {
    panic(err)
  }
  return ToJSON
}

// FromJSON used to unmarshalling Book type
func FromJSON(data []byte) Book {
  b := Book{}
  err := json.Unmarshal(data, &b)
  if err != nil {
    panic(err)
  }
  return b
}

// BooksHandleFunc is a function that handle /api/books
func BooksHandleFunc(w http.ResponseWriter, r *http.Request) {
  switch method := r.Method; method {
  case http.MethodGet:
    books := AllBooks()
    writeJSON(w, books)
  case http.MethodPost:
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
    }
    book := FromJSON(body)
    isbn, created := CreateBook(book)
    if created {
      w.Header().Add("Location", "api/books/"+isbn)
      w.WriteHeader(http.StatusCreated)
    } else {
      w.WriteHeader(http.StatusConflict)
    }
  default:
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("Unsupported request method"))
  }
}

// BookHandleFunc handle /api/books/<isbn>
func BookHandleFunc(w http.ResponseWriter, r *http.Request) {
  isbn := r.URL.Path[len("/api/books/"):]

  switch method := r.Method; method {
  case http.MethodGet:
    book, found := GetBook(isbn)
    if found {
      writeJSON(w, book)
    } else {
      w.WriteHeader(http.StatusNotFound)
    }
  case http.MethodPut:
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
    }
    book := FromJSON(body)
    exists := UpdateBook(isbn, book)
    if exists {
      w.WriteHeader(http.StatusOK)
    } else {
      w.WriteHeader(http.StatusNotFound)
    }
  case http.MethodDelete:
    DeleteBook(isbn)
    w.WriteHeader(http.StatusOK)
  default:
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("unsupported request method"))
  }
}

// AllBooks return a slice of all books
func AllBooks() []Book {
  values := make([]Book, len(books))
  idx := 0
  for _, book := range books {
    values[idx] = book
    idx++
  }

  return values
}

func writeJSON(w http.ResponseWriter, i interface{}) {
  b, err := json.Marshal(i)
  if err != nil {
    panic(err)
  }
  w.Header().Add("Content-Type", "application/json;charset=utf-8")
  w.Write(b)
}

// GetBook return the book fo a given isbn
func GetBook(isbn string) (Book, bool) {
  book, found := books[isbn]
  return book, found
}

// CreateBook creates a new book if it doesn't exist
func CreateBook(book Book) (string, bool) {
  _, exists := books[book.ISBN]
  if exists {
    return "", false
  }
  books[book.ISBN] = book
  return book.ISBN, true
}

// UpdateBook updates an existing book
func UpdateBook(isbn string, book Book) bool {
  _, exists := books[isbn]
  if exists {
    books[isbn] = book
  }
  return exists
}

// DeleteBook delete books from a map by isbn key
func DeleteBook(isbn string) {
  delete(books, isbn)
}
