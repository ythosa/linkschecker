package links

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

// Crawl ...
func Crawl(url string) []string {
	//fmt.Println(url)

	list, err := Extract(url)
	if err != nil {
		log.Print(err)
	}

	return list
}

func Extract(url string) ([]string, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Link with URL {%s} responds status code: %d\n", url, response.StatusCode)
	}

	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Println(err.Error())
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
					fmt.Printf("Error: %s \n", err.Error())
					continue
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
