package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/html"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func getLinks(url string) ([]string, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	var extractLinks func(*html.Node)
	extractLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					absoluteURL := resolveURL(attr.Val, url)
					links = append(links, absoluteURL)
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractLinks(c)
		}
	}
	extractLinks(doc)

	return links, nil
}

func resolveURL(href, baseURL string) string {
	parsedBase, err := url.Parse(baseURL)
	if err != nil {
		return href
	}

	parsedHref, err := url.Parse(href)
	if err != nil {
		return href
	}

	return parsedBase.ResolveReference(parsedHref).String()
}

func getStatusCode(url string) (int, string) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "Error en solicitud"
	}

	resp, err := client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			return 0, "Timeout"
		}
		return 0, "Error"
	}
	defer resp.Body.Close()
	return resp.StatusCode, http.StatusText(resp.StatusCode)
}

func main() {
	var url string
	fmt.Print("Ingrese su enlace a screapear: ")
	fmt.Scanf("%s", &url)
	links, err := getLinks(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Enlaces encontrados:")
	for _, link := range links {

		statusCode, statusText := getStatusCode(link)
		// if err != nil {
		// 	fmt.Println("Error: ", err)
		// 	return
		// }

		//if statusCode != http.StatusOK {
		fmt.Println(link, "Code:", statusCode, statusText)
		//}
	}

	fmt.Println("Scrapeo terminado.")
}
