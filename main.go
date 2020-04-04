package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aiomonitors/godiscord"
	proxymanager "github.com/aiomonitors/goproxymanager"
	"github.com/fatih/color"
)

type Config struct {
	Webhooks  []godiscord.Webhook `json:"webhooks"`
	Links     []string            `json:"links"`
	Color     string              `json:"color"`
	Groupname string              `json:"groupname"`
}

type Monitor struct {
	Config       Config
	Availability map[string]bool
	Manager      proxymanager.ProxyManager
	UseProxies   bool
}

type ProductInfo struct {
	Name         string `json:"name"`
	SKU          string `json:"sku"`
	Availability bool   `json:"availability"`
	Link         string `json:"link"`
	Image        string `json:"image"`
	Exec         string `json:"exec"`
}

func getTime() string {
	return time.Now().Format("15:04:05")
}

func printErr(err error) {
	color.Red("[ %s ] Error in monitor %v", getTime(), err)
}

func greenMessage(msg string) {
	color.Green("[ %s ] %s", getTime(), msg)
}

func yellowMessage(msg string) {
	color.Yellow("[ %s ] %s", getTime(), msg)
}

func redMessage(msg string) {
	color.Red("[ %s ] %s", getTime(), msg)
}

func availability(link string, m Monitor) (*ProductInfo, error) {
	start := time.Now()

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
	//Sets up http client with proxies if monitor uses proxies
	var client *http.Client
	if m.UseProxies == true {
		proxy, proxyErr := m.Manager.NextProxy()
		if proxyErr != nil {
			return &ProductInfo{}, proxyErr
		}
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return &ProductInfo{}, err
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	} else {
		client = &http.Client{}
	}

	//Sets up HTTP requests
	req, reqErr := http.NewRequest("GET", link, nil)
	if reqErr != nil {
		return &ProductInfo{}, reqErr
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	//Does request
	res, resErr := client.Do(req)
	if resErr != nil {
		return &ProductInfo{}, resErr
	}
	defer res.Body.Close()

	document, documentErr := goquery.NewDocumentFromReader(res.Body)
	if documentErr != nil {
		return &ProductInfo{}, documentErr
	}

	name := document.Find("div.sku-title > h1").Text()
	sku := document.Find("div.sku.product-data > span.product-data-value.body-copy").Text()
	buttonText := document.Find(".fulfillment-add-to-cart-button > div").Find("button").Text()
	image, _ := document.Find("div.primary-image-container > button > img").Attr("src")

	elapsed := fmt.Sprintf("%s", time.Since(start))

	if buttonText != "" && buttonText == "Add to Cart" {
		return &ProductInfo{name, sku, true, link, image, elapsed}, nil
	}
	return &ProductInfo{name, sku, false, link, image, elapsed}, nil
}

func NewMonitor(pathToConfig string, proxyPath string) (*Monitor, error) {
	m := Monitor{}

	//Proxy initialization
	if proxyPath != "" {
		manager, managerErr := proxymanager.NewManager(proxyPath)
		if managerErr != nil {
			m.Manager = *manager
		} else {
			m.UseProxies = false
			color.Red("Error loading proxy file")
		}
	} else {
		m.UseProxies = false
	}
	//Config initialization
	if pathToConfig != "" {
		file, openErr := ioutil.ReadFile(pathToConfig)
		if openErr != nil {
			return nil, openErr
		}
		var c Config
		unmarshalError := json.Unmarshal(file, &c)
		if unmarshalError != nil {
			return nil, unmarshalError
		}
		m.Config = c
	} else {
		return nil, errors.New("Need to provide a config file")
	}
	//Initialize links
	m.Availability = make(map[string]bool)
	for _, link := range m.Config.Links {
		req, reqErr := availability(link, m)
		if reqErr != nil {
			redMessage(fmt.Sprintf("%s error getting link", link))
			m.Availability[link] = true
		} else {
			greenMessage(fmt.Sprintf("Initialized %s Availability: %v", req.Name, req.Availability))
			m.Availability[link] = req.Availability
		}
	}
	return &m, nil
}

func (m Monitor) Monitor() {
	i := true
	for i == true {
		for _, link := range m.Config.Links {
			req, reqErr := availability(link, m)
			yellowMessage(fmt.Sprintf("[ %s ] Monitoring %s", req.Exec, req.SKU))
			if reqErr != nil {
				redMessage(fmt.Sprintf("Error initializing link %s %v", link, reqErr))
			} else {
				if m.Availability[link] == false && req.Availability == true {
					m.Availability[link] = true
					greenMessage(fmt.Sprintf("%s in stock!", req.SKU))
					m.SendEmbed(req)
				} else if m.Availability[link] == true && req.Availability == false {
					m.Availability[link] = false
					redMessage(fmt.Sprintf("%s is out of stock", req.SKU))
				}
			}
		}
		time.Sleep(time.Millisecond * 1500)
	}
}

func (m Monitor) SendEmbed(p *ProductInfo) {
	for _, webhook := range m.Config.Webhooks {
		e := godiscord.NewEmbed(p.Name, fmt.Sprintf("In Stock: True"), p.Link)
		e.SetFooter(webhook.Text, webhook.IconURL)
		e.SetColor("#16A085")
		e.SetAuthor("Best Buy Monitor", "", "")
		e.SetThumbnail(p.Image)
		e.SendToWebhook(webhook.URL)
	}

}

func main() {
	var configFile = flag.String("config", "", "link to check")
	var proxyFile = flag.String("proxy", "", "proxy file")

	flag.Parse()
	m, mErr := NewMonitor(*configFile, *proxyFile)
	if mErr != nil {
		panic(mErr)
	}
	m.Monitor()
}
