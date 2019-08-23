package pkg

import (
	"errors"
	"strings"
	"sync"
)

type mutexedID struct {
	id    int
	mutex sync.Mutex
}

var globalID mutexedID

type mutexedURLs struct {
	shortened map[string]string
	mutex     sync.Mutex
}

var globalURLs mutexedURLs

func init() {
	globalURLs.shortened = make(map[string]string)
}

func getNextID() int {
	globalID.mutex.Lock()
	id := globalID.id
	globalID.id++
	globalID.mutex.Unlock()
	return id
}

func addGlobalURL(shortened, url string) {
	globalURLs.mutex.Lock()
	globalURLs.shortened[shortened] = url
	globalURLs.mutex.Unlock()
}

var chars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type Service interface {
	createShortenedURL(url string) (string, error)
	retrieveURL(shortened string) (string, error)
}

type ServiceImpl struct{}

func NewServiceImpl() *ServiceImpl {
	return &ServiceImpl{}
}

func (service *ServiceImpl) createShortenedURL(url string) (string, error) {
	if url == "" && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", errors.New("incorrent url string")
	}
	n := getNextID()
	var builder strings.Builder
	builder.WriteString("/url/shortened/")
	charLen := len(chars)
	for {
		builder.WriteByte(chars[n%charLen])
		n = n / charLen
		if n <= 0 {
			break
		}
	}
	shortened := builder.String()
	addGlobalURL(shortened, url)
	return shortened, nil
}

func (service *ServiceImpl) retrieveURL(shortened string) (string, error) {
	if url, ok := globalURLs.shortened[shortened]; !ok {
		return "", errors.New("shortened url doesn't exist, shortened:" + shortened)
	} else {
		return url, nil
	}
}
