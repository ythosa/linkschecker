package links

import (
    "net/http"

    "golang.org/x/net/html"
)

// CheckURL checks passed url in arguments and returns specified error if it is incorrect.
func CheckURL(url ParsingURL) (*http.Response, *html.Node, error) {
    response, err := http.Get(string(url))
    if err != nil {
        return nil, nil, NewUnreachableSiteException(string(url))
    }

    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return nil, nil, NewBadStatusCodeException(string(url), response.StatusCode)
    }

    doc, err := html.Parse(response.Body)
    if err != nil {
        return nil, nil, NewInvalidResponseTypeException(string(url))
    }

    return response, doc, nil
}

// Extract extracts all links from passed URL.
func Extract(response *http.Response, doc *html.Node) []ParsingURL {
    var links []ParsingURL
    visitNode := func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, a := range n.Attr {
                if a.Key != "href" {
                    continue
                }

                link, err := response.Request.URL.Parse(a.Val)
                if err != nil {
                    return
                }

                links = append(links, ParsingURL(link.String()))
            }
        }
    }

    ForEachNode(doc, visitNode)

    return links
}

// ForEachNode recursive runs throw DOM tree and finds links.
func ForEachNode(node *html.Node, f func(n *html.Node)) {
    if f != nil {
        f(node)
    }

    for c := node.FirstChild; c != nil; c = c.NextSibling {
        ForEachNode(c, f)
    }
}
