package pkg

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var service = NewServiceImpl()

func TestShortenIncorrectURL(t *testing.T) {
	incorrectURL := ""
	_, err := service.createShortenedURL(incorrectURL)
	assert.NotNil(t, err)
}

func TestShortenedURL(t *testing.T) {
	urls := []string{
		"https://www.google.com/imgres?imgurl=https%3A%2F%2Fcdn0.tnwcdn.com%2Fwp-content%2Fblogs.dir%2F1%2Ffiles%2F2018%2F07%2Fgo.png&imgrefurl=https%3A%2F%2Fthenextweb.com%2Fdd%2F2018%2F07%2F25%2Fgoogles-new-package-helps-developers-use-golang-in-cloud-apps%2F&docid=oxOLgamvbK7skM&tbnid=xdNCv_8YWeIglM%3A&vet=10ahUKEwjrv9uN1ZjkAhUL2aYKHTuKDHkQMwh7KAEwAQ..i&w=607&h=318&bih=860&biw=1680&q=golang&ved=0ahUKEwjrv9uN1ZjkAhUL2aYKHTuKDHkQMwh7KAEwAQ&iact=mrc&uact=8",
		"https://www.google.com/imgres?imgurl=https%3A%2F%2Fcamo.githubusercontent.com%2F0e3c4976eb4b8ec80e285608a7f23744408a0ffb%2F68747470733a2f2f736563757265676f2e696f2f696d672f676f7365632e706e67&imgrefurl=https%3A%2F%2Fgithub.com%2Fsecurego%2Fgosec&docid=yDGeIwH1ck8kqM&tbnid=zisxOtuU0oP3BM%3A&vet=10ahUKEwjrv9uN1ZjkAhUL2aYKHTuKDHkQMwiMASgRMBE..i&w=2362&h=2362&bih=860&biw=1680&q=golang&ved=0ahUKEwjrv9uN1ZjkAhUL2aYKHTuKDHkQMwiMASgRMBE&iact=mrc&uact=8",
		"https://www.google.com/imgres?imgurl=https%3A%2F%2Fpbs.twimg.com%2Fmedia%2FDfQF_XOX4AAjzxW.png%3Alarge&imgrefurl=https%3A%2F%2Ftwitter.com%2Fjedisct1%2Fstatus%2F1006057491237687296&docid=r9Aei7-vvv2xCM&tbnid=8vJ8kR00l2cFMM%3A&vet=1&w=686&h=600&bih=860&biw=1680&ved=2ahUKEwjH65Wp1ZjkAhUIpIsKHct3CpsQxiAoBHoECAEQHw&iact=c&ictx=1",
	}
	shortenedUrls := make([]string, len(urls))
	for i, url := range urls {
		shortened, err := service.createShortenedURL(url)
		assert.Nil(t, err)
		shortenedUrls[i] = shortened

	}
	url, error := service.retrieveURL(shortenedUrls[0])
	assert.Nil(t, error)
	assert.Equal(t, urls[0], url)
}
