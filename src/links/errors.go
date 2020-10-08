package links

import "fmt"

type badStatusCodeException struct {
    url        ParsingURL
    statusCode int
}

type unreachableSiteException struct {
    url ParsingURL
}

type invalidResponseTypeException struct {
    url ParsingURL
}

func NewBadStatusCodeException(url ParsingURL, statusCode int) error {
    return &badStatusCodeException{url: url, statusCode: statusCode}
}

func NewUnreachableSiteException(url ParsingURL) error {
    return &unreachableSiteException{url: url}
}

func NewInvalidResponseTypeException(url ParsingURL) error {
    return &invalidResponseTypeException{url: url}
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
