package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type ServiceStub struct{}

func (stub *ServiceStub) createShortenedURL(url string) (string, error) {
	if url == "" {
		return "", errors.New("error")
	}
	return "test", nil
}

func (stub *ServiceStub) retrieveURL(shortened string) (string, error) {
	if shortened == "" {
		return "", errors.New("error")
	} else if shortened == "/url/shortened/" {
		return "", errors.New("error")
	}
	return "test", nil
}

var handler = NewHandlerImpl(&ServiceStub{})

func TestHandlerGetShortenedURLWrongMethod(t *testing.T) {
	request := httptest.NewRequest("GET", "/url/shortened", nil)
	writer := httptest.NewRecorder()
	handler.GetShortenedURL(writer, request)
	response := writer.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestHandlerGetShortenedURLBadRequestError(t *testing.T) {
	request := httptest.NewRequest("POST", "/url/shortened", strings.NewReader("test"))
	writer := httptest.NewRecorder()
	handler.GetShortenedURL(writer, request)
	response := writer.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestHandlerGetShortenedURLInternalServerError(t *testing.T) {
	jsonBytes, _ := json.Marshal(FullURL{})
	request := httptest.NewRequest("POST", "/url/shortened", bytes.NewReader(jsonBytes))
	writer := httptest.NewRecorder()
	handler.GetShortenedURL(writer, request)
	response := writer.Result()
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestHandlerGetShortenedURLStatusOK(t *testing.T) {
	jsonBytes, _ := json.Marshal(FullURL{URL: "http://"})
	request := httptest.NewRequest("POST", "/url/shortened", bytes.NewReader(jsonBytes))
	writer := httptest.NewRecorder()
	handler.GetShortenedURL(writer, request)
	response := writer.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandlerGetFullURLWrongMethod(t *testing.T) {
	request := httptest.NewRequest("GET", "/url/full", nil)
	writer := httptest.NewRecorder()
	handler.GetFullURL(writer, request)
	response := writer.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestHandlerGetFullURLBadRequestError(t *testing.T) {
	request := httptest.NewRequest("POST", "/url/full", strings.NewReader("test"))
	writer := httptest.NewRecorder()
	handler.GetFullURL(writer, request)
	response := writer.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestHandlerGetFullURLInternalServerError(t *testing.T) {
	jsonBytes, _ := json.Marshal(ShortenedURL{})
	request := httptest.NewRequest("POST", "/url/full", bytes.NewReader(jsonBytes))
	writer := httptest.NewRecorder()
	handler.GetFullURL(writer, request)
	response := writer.Result()
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestHandlerGetFullURLStatusOK(t *testing.T) {
	jsonBytes, _ := json.Marshal(ShortenedURL{Shortened: "http://"})
	request := httptest.NewRequest("POST", "/url/full", bytes.NewReader(jsonBytes))
	writer := httptest.NewRecorder()
	handler.GetFullURL(writer, request)
	response := writer.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandlerGetRedirectURLWrongMethod(t *testing.T) {
	request := httptest.NewRequest("PUT", "/url/shortened/", nil)
	writer := httptest.NewRecorder()
	handler.GetRedirectURL(writer, request)
	response := writer.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestHandlerGetRedirectURLNotFound(t *testing.T) {
	request := httptest.NewRequest("GET", "/url/shortened/", nil)
	writer := httptest.NewRecorder()
	handler.GetRedirectURL(writer, request)
	response := writer.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestHandlerGetRedirectURLStatusFound(t *testing.T) {
	request := httptest.NewRequest("GET", "/url/shortened/a", nil)
	writer := httptest.NewRecorder()
	handler.GetRedirectURL(writer, request)
	response := writer.Result()
	assert.Equal(t, http.StatusFound, response.StatusCode)
}
