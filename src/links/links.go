package links

import (
    "net/http"

    "golang.org/x/net/html"
)

func Check(url string) (*http.Response, *html.Node, error) {
    response, err := http.Get(url)
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

func Extract(url string) ([]string, error) {
    response, doc, err := Check(url)
    if err != nil {
        return nil, err
    }

    var links []string

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

                links = append(links, link.String())
            }
        }
    }

    ForEachNode(doc, visitNode)

    return links, nil
}

func ForEachNode(node *html.Node, f func(n *html.Node)) {
    if f != nil {
        f(node)
    }

    for c := node.FirstChild; c != nil; c = c.NextSibling {
        ForEachNode(c, f)
    }
}
