package links

import (
    "golang.org/x/net/html"
    "net/http"
    "strings"
    "sync"
)

// CheckURL checks passed url in arguments and returns specified error if it is incorrect.
func CheckURL(url ParsingURL) (*http.Response, *html.Node, error) {
    response, err := http.Get(string(url))
    if err != nil {
        return nil, nil, NewUnreachableSiteException(url)
    }

    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return nil, nil, NewBadStatusCodeException(url, response.StatusCode)
    }

    doc, err := html.Parse(response.Body)
    if err != nil {
        return nil, nil, NewInvalidResponseTypeException(url)
    }

    return response, doc, nil
}

// FindBrokenLinks finds broken links on site and returns array of type BadURL as result.
func FindBrokenLinks(baseURL ParsingURL) []BadURL {
    worklist := make(chan []ParsingURL)
    errlist := make([]BadURL, 0)

    var n int

    n++
    go func() {
        worklist <- []ParsingURL{baseURL}
    }()

    seen := make(map[ParsingURL]bool)
    mux := sync.Mutex{}
    for ; n > 0; n-- {
        slist := <-worklist
        if slist == nil {
            continue
        }

        for _, link := range slist {
            if !seen[link] {
                seen[link] = true
                n++
                go func(link ParsingURL) {
                    res, doc, err := CheckURL(link)
                    if err != nil {
                        mux.Lock()
                        errlist = append(errlist, BadURL{ParsingURL: link, Err: err})
                        mux.Unlock()

                        worklist <- nil
                        return
                    }

                    if strings.HasPrefix(string(link), string(baseURL)) {
                        worklist <- Extract(res, doc)
                    } else {
                        worklist <- nil
                    }
                }(link)
            }
        }
    }

    return errlist
}
