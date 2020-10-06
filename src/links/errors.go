package links

import "fmt"

type badStatusCodeException struct {
    url        string
    statusCode int
}

type unreachableSiteException struct {
    url string
}

type invalidResponseTypeException struct {
    url string
}

type invalidURLException struct {
    url string
}

func NewBadStatusCodeException(url string, statusCode int) error {
    return &badStatusCodeException{url: url, statusCode: statusCode}
}

func NewUnreachableSiteException(url string) error {
    return &unreachableSiteException{url: url}
}

func NewInvalidResponseTypeException(url string) error {
    return &invalidResponseTypeException{url: url}
}

func NewInvalidURLException(url string) error {
    return &invalidURLException{url: url}
}

func (e *badStatusCodeException) Error() string {
    return fmt.Sprintf("%s - bad status code response - %d", e.url, e.statusCode)
}

func (e *unreachableSiteException) Error() string {
    return fmt.Sprintf("%s - is unreachable", e.url)
}

func (e *invalidResponseTypeException) Error() string {
    return fmt.Sprintf("%s - invalid response body", e.url)
}

func (e *invalidURLException) Error() string {
    return fmt.Sprintf("%s - invalid format", e.url)
}
