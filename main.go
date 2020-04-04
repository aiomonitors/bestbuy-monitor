package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

func availability(link string) (bool, error) {
	headers := map[string]string{
		"authority":                 "www.bestbuy.com",
		"cache-control":             "max-age=0",
		"upgrade-insecure-requests": "1",
		"user-agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
		"sec-fetch-dest":            "document",
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"sec-fetch-site":            "same-origin",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-user":            "?1",
		"referer":                   "https://www.bestbuy.com/site/searchpage.jsp?st=nintendo+switch&_dyncharset=UTF-8&id=pcat17071&type=page&sc=Global&cp=1&nrp=&sp=&qp=&list=n&af=true&iht=y&usc=All+Categories&ks=960&keys=keys",
		"accept-language":           "en-US,en;q=0.9",
	}
	req, reqErr := http.NewRequest("GET", link, nil)
	if reqErr != nil {
		return false, reqErr
	}

	client := &http.Client{}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, resErr := client.Do(req)
	if resErr != nil {
		return false, resErr
	}
	defer res.Body.Close()

	document, documentErr := goquery.NewDocumentFromReader(res.Body)
	if documentErr != nil {
		return false, documentErr
	}

	name := document.Find("div.sku-title > h1").Text()
	sku := document.Find("div.sku.product-data > span.product-data-value.body-copy").Text()
	buttonText := document.Find(".fulfillment-add-to-cart-button > div").Find("button").Text()

	color.Blue("Name: %s", name)
	color.Blue("SKU: %s", sku)
	color.Blue("Button Text: %s", buttonText)
	if buttonText != "" && buttonText == "Add to Cart" {
		return true, nil
	}
	return false, nil //
}

func main() {
	var link = flag.String("link", "", "link to check")
	flag.Parse()

	if *link == "" {
		panic(errors.New("missing link"))
	} else {
		availability, err := availability(*link)
		if err != nil {
			panic(err)
		}
		fmt.Println(availability)
	}
}
