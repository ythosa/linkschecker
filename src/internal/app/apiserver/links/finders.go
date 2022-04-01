package links

import (
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

// isUnderBaseURL returns true if link placed under baseURL, otherwise - false.
func isUnderBaseURL(link ParsingURL, baseURL ParsingURL) bool {
	return strings.HasPrefix(string(link), string(baseURL))
}

// CheckURL checks passed url in arguments and returns specified error if it is incorrect.
func CheckURL(url ParsingURL) (*http.Response, *html.Node, error) {
	response, err := http.Get(string(url))
	if err != nil {
		return nil, nil, NewUnreachableSiteError(url)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, nil, NewBadStatusCodeError(url, response.StatusCode)
	}

	doc, err := html.Parse(response.Body)
	if err != nil {
		return nil, nil, NewInvalidResponseTypeError(url)
	}

	return response, doc, nil
}

// FindBrokenLinks finds broken links on site and returns array of type BrokenURL as result.
func FindBrokenLinks(baseURL ParsingURL) []BrokenURL {
	workList := make(chan []ParsingURL) // channel of links to be checked
	errList := make([]BrokenURL, 0)     // slice of broken links

	var workListElementsNum int
	var mux sync.Mutex

	workListElementsNum++
	go func() {
		workList <- []ParsingURL{baseURL} // Start with baseURL
	}()

	seenURLs := make(map[ParsingURL]bool) // map of seen URLs
	for ; workListElementsNum > 0; workListElementsNum-- {
		extractedLinks := <-workList
		if extractedLinks == nil {
			continue
		}

		for _, link := range extractedLinks {
			if !seenURLs[link] {
				seenURLs[link] = true
				workListElementsNum++

				go func(link ParsingURL) {
					res, doc, err := CheckURL(link)
					if err != nil {
						mux.Lock()
						errList = append(errList, BrokenURL{ParsingURL: link, Error: err})
						mux.Unlock()
						workList <- nil

						return
					}

					if isUnderBaseURL(link, baseURL) {
						workList <- Extract(res, doc)

						return
					}

					workList <- nil
				}(link)
			}
		}
	}

	return errList
}
