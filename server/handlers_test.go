package server

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"root/db"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func getFile(is *is.I, fileName string) *os.File {
	f, err := os.Open(fileName)
	if err != nil {
		is.NoErr(err) // cannot open file
		is.Fail()
	}
	return f
}

func prepareForm(is *is.I, data map[string]io.Reader) (*bytes.Buffer, *multipart.Writer) {

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	defer w.Close()

	for k, r := range data {

		var fw io.Writer
		var err error

		// Try to defer-close the reader.
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}

		// If the reader is a os.File, add it as one.
		if x, ok := r.(*os.File); ok {
			fw, err = w.CreateFormFile(k, x.Name())
			if err != nil {
				is.NoErr(err) // cannot create form file
				is.Fail()
			}
		} else {
			// Add other fields.
			fw, err = w.CreateFormField(k)
			if err != nil {
				is.NoErr(err) // cannot create form file
				is.Fail()
			}
		}

		if _, err = io.Copy(fw, r); err != nil {
			is.NoErr(err) // cannot copy
			is.Fail()
		}
	}

	return &b, w
}

func TestPostsCreate(t *testing.T) {
	is := is.New(t)

	// Create the server with a mock db.
	db := db.NewMockDB()
	s := NewServer(db)

	// Create a test table.
	tt := []struct {
		formData map[string]io.Reader
		expected int
	}{
		{
			map[string]io.Reader{
				"title":       strings.NewReader("Hello, World!"),
				"description": strings.NewReader("Hello, World!"),
				"thumbnail":   getFile(is, "for_handlers_test.png"),
			},
			http.StatusCreated,
		},
		{
			map[string]io.Reader{
				"title":     strings.NewReader("Hello, World!"),
				"thumbnail": getFile(is, "for_handlers_test.png"),
			},
			http.StatusUnprocessableEntity,
		},
		{
			map[string]io.Reader{
				"description": strings.NewReader("Hello, World!"),
				"thumbnail":   getFile(is, "for_handlers_test.png"),
			},
			http.StatusUnprocessableEntity,
		},
		{
			map[string]io.Reader{
				"title":       strings.NewReader("Hello, World!"),
				"description": strings.NewReader("Hello, World!"),
			},
			http.StatusUnprocessableEntity,
		},
	}

	// Run the test cases.
	for _, tc := range tt {

		// Prepare the form.
		b, mw := prepareForm(is, tc.formData)

		// Create the request and response-recorder.
		r := httptest.NewRequest("POST", "/posts", b)
		r.Header.Set("Content-Type", mw.FormDataContentType())
		w := httptest.NewRecorder()

		// Do the request.
		s.ServeHTTP(w, r)

		// Run tests.
		is.Equal(w.Code, tc.expected)
	}
}

func TestPostsGet(t *testing.T) {
	is := is.New(t)

	// Create the server with a mock db.
	db := db.NewMockDB()
	s := NewServer(db)

	// Create a test table.
	tt := []struct {
		postID   string
		expected int
	}{
		{
			"dosenotexist",
			http.StatusNotFound,
		},
		{
			"exists",
			http.StatusOK,
		},
	}

	// Run the test cases.
	for _, tc := range tt {

		// Create the request and response-recorder.
		r := httptest.NewRequest("GET", "/posts/"+tc.postID, nil)
		w := httptest.NewRecorder()

		// Do the request.
		s.ServeHTTP(w, r)

		// Run tests.
		is.Equal(w.Code, tc.expected)
	}
}

func TestPostsGetAll(t *testing.T) {
	is := is.New(t)

	// Create the server with a mock db.
	db := db.NewMockDB()
	s := NewServer(db)

	// Create a test table.
	tt := []struct {
		offset   string
		limit    string
		expected int
	}{
		{
			offset:   "0",
			limit:    "10",
			expected: http.StatusOK,
		},
		{
			offset:   "10000000000000000000000000000000",
			limit:    "100000000000000000000000000000000000000000000000000",
			expected: http.StatusUnprocessableEntity,
		},
	}

	// Run the test cases.
	for _, tc := range tt {

		// Create the request and response-recorder.
		url := fmt.Sprintf("/posts/%s/%s", tc.offset, tc.limit)
		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		// Do the request.
		s.ServeHTTP(w, r)

		// Run tests.
		is.Equal(w.Code, tc.expected)
	}
}

func TestPostsDelete(t *testing.T) {
	is := is.New(t)

	// Create the server with a mock db.
	db := db.NewMockDB()
	s := NewServer(db)

	// Create a test table.
	tt := []string{"exists", "dosenotexist"}

	// Create the request and response-recorder.
	for _, tc := range tt {
		r := httptest.NewRequest("DELETE", "/posts/"+tc, nil)
		w := httptest.NewRecorder()

		// Do the request.
		s.ServeHTTP(w, r)

		// Run tests.
		is.Equal(w.Code, http.StatusNoContent)
	}
}
