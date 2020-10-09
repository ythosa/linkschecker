package links

import "fmt"

// badStatusCodeException is type of error where the server under URL
// answers with status code other than 200.
type badStatusCodeException struct {
    url        ParsingURL
    statusCode int
}

// unreachableSiteException is type of error where the server under URL is unreachable.
type unreachableSiteException struct {
    url ParsingURL
}

// invalidResponseTypeException is type of error where the server under URL
// returns an invalid response body.
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
