package links

import "fmt"

// badStatusCodeError is type of error where the server under URL
// answers with status code other than 200.
type badStatusCodeError struct {
	url        ParsingURL
	statusCode int
}

// unreachableSiteError is type of error where the server under URL is unreachable.
type unreachableSiteError struct {
	url ParsingURL
}

// invalidResponseTypeError is type of error where the server under URL
// returns an invalid response body.
type invalidResponseTypeError struct {
	url ParsingURL
}

func NewBadStatusCodeError(url ParsingURL, statusCode int) error {
	return &badStatusCodeError{url: url, statusCode: statusCode}
}

func NewUnreachableSiteError(url ParsingURL) error {
	return &unreachableSiteError{url: url}
}

func NewInvalidResponseTypeError(url ParsingURL) error {
	return &invalidResponseTypeError{url: url}
}

func (e *badStatusCodeError) Error() string {
	return fmt.Sprintf("%s - bad status code response - %d", e.url, e.statusCode)
}

func (e *unreachableSiteError) Error() string {
	return fmt.Sprintf("%s - is unreachable", e.url)
}

func (e *invalidResponseTypeError) Error() string {
	return fmt.Sprintf("%s - invalid response body", e.url)
}
