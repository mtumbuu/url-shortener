package pkg

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Handler interface {
	GetShortenedURL(http.ResponseWriter, *http.Request)
	GetFullURL(http.ResponseWriter, *http.Request)
	GetRedirectURL(http.ResponseWriter, *http.Request)
}

type HandlerImpl struct {
	Service Service
}

func NewHandlerImpl(service Service) *HandlerImpl {
	return &HandlerImpl{Service: service}
}

func (handler *HandlerImpl) GetShortenedURL(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Redirect(writer, request, "/notFound", http.StatusNotFound)
	}
	decoder := json.NewDecoder(request.Body)
	urlRequest := &FullURL{}
	if err := decoder.Decode(urlRequest); err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(writer)
		encoder.Encode("bad request")
		return
	}
	if shortened, err := handler.Service.createShortenedURL(urlRequest.URL); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(writer)
		encoder.Encode("something went wrong")
	} else {
		encoder := json.NewEncoder(writer)
		encoder.Encode(&ShortenedURL{Shortened: request.Host + shortened})
	}
}

func (handler *HandlerImpl) GetFullURL(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Redirect(writer, request, "/notFound", http.StatusNotFound)
	}
	decoder := json.NewDecoder(request.Body)
	urlRequest := &ShortenedURL{}
	if err := decoder.Decode(urlRequest); err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(writer)
		encoder.Encode("bad request")
		return
	}
	stripped := strings.TrimPrefix(urlRequest.Shortened, request.Host)
	if url, err := handler.Service.retrieveURL(stripped); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(writer)
		encoder.Encode("not found")
	} else {
		encoder := json.NewEncoder(writer)
		encoder.Encode(&FullURL{URL: url})
	}
}

func (handler *HandlerImpl) GetRedirectURL(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Redirect(writer, request, "/notFound", http.StatusNotFound)
	}
	if url, err := handler.Service.retrieveURL(request.URL.String()); err != nil {
		http.Redirect(writer, request, "/notFound", http.StatusNotFound)
	} else {
		http.Redirect(writer, request, url, http.StatusFound)
	}
}
