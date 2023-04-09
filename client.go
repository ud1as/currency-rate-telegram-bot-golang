package main

import (
	"net/http"
)

// fetchFunction is a function that mimics http.Get() method
type fetchFunction func(url string) (resp *http.Response, err error)

type Client interface {
	GetRate(code string) (float32, float32, error)
	SetFetchFunction(fetchFunction)
}

type client struct {
	fetch fetchFunction
}

func (s client) GetRate(code string) (float32, float32, error) {
	rate, difference, err := getRate(code, s.fetch)
	if err != nil {
		return 0, 0, err
	}
	return rate, difference, nil
}

func (s client) SetFetchFunction(f fetchFunction) {
	s.fetch = f
}

// NewClient creates a new rates service instance
func NewClient() Client {
	return client{http.Get}
}
