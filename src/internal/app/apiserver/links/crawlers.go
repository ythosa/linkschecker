package links

import (
    "net/http"

    "golang.org/x/net/html"
)

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
