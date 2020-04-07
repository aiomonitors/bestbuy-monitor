package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	bbresponse "github.com/aiomonitors/bestbuy-monitor/ca/models"
	"github.com/aiomonitors/godiscord"
	proxymanager "github.com/aiomonitors/goproxymanager"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
)

var headers = map[string]string{
	"authority":        "www.bestbuy.ca",
	"accept":           "application/json, text/javascript, */*; q=0.01",
	"sec-fetch-dest":   "empty",
	"x-requested-with": "XMLHttpRequest",
	"user-agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
	"sec-fetch-site":   "same-origin",
	"sec-fetch-mode":   "cors",
	"referer":          "https://www.bestbuy.ca/api/v2/json/search?categoryid=&currentRegion=ON&include=facets%2C%20redirects&lang=en-CA&page=1&pageSize=24&path=&query=13817625&exp=&sortBy=relevance&sortDir=desc",
	"accept-language":  "en-US,en;q=0.9",
}

var headers2 = map[string]string{
	"authority":       "www.bestbuy.ca",
	"user-agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
	"accept":          "application/json, text/javascript, */*; q=0.01",
	"referer":         "https://www.bestbuy.ca/en-ca/product/nintendo-switch-console-with-grey-joy-con/13817626",
	"accept-language": "en-US,en;q=0.9",
	"Cache-Control":   "no-cache",
}

type ProductInfo struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Price string `json:"price"`
	SKU   string `json:"sku"`
	Link  string `json:"link"`
	Exec  string `json:"exec"`
}

type ProductStock struct {
	Available bool   `json:"available"`
	Stock     int    `json:"stock"`
	Exec      string `json:"exec"`
}

type Config struct {
	Webhooks []godiscord.Webhook `json:"webhooks"`
	SKUs     []string            `json:"skus"`
	Color    string              `json:"color"`
	Links    []string            `json:"links"`
}

type Monitor struct {
	Config       Config                    `json:"config"`
	Availability map[string]bool           `json:"availability`
	ProductInfo  map[string]ProductInfo    `json:"product_info"`
	Manager      proxymanager.ProxyManager `json:"manager"`
	UseProxies   bool                      `json:"use_proxies"`
	Client       http.Client               `json:"client"`
}

type MonitorInterface interface {
	GetInfo() (interface{}, error)
	UpdateClient() error
}

func GetTime() string {
	return time.Now().Format("15:04:05")
}

func PrintErr(err error) {
	color.Red("[ %s ] Error in monitor %v", GetTime(), err)
}

func GreenMessage(msg string) {
	color.Green("[ %s ] %s", GetTime(), msg)
}

func YellowMessage(msg string) {
	color.Yellow("[ %s ] %s", GetTime(), msg)
}

func RedMessage(msg string) {
	color.Red("[ %s ] %s", GetTime(), msg)
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
	m.UpdateClient()
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
	return &m, nil
}

func (m *Monitor) UpdateClient() error {
	if m.UseProxies == true {
		proxy, proxyErr := m.Manager.NextProxy()
		if proxyErr != nil {
			return proxyErr
		}
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return err
		}
		m.Client = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	} else {
		m.Client = http.Client{}
	}
	return nil
}

func (m Monitor) GetInfo(sku string) (*ProductInfo, error) {
	start := time.Now()
	m.UpdateClient()
	link := fmt.Sprintf("https://www.bestbuy.ca/api/v2/json/search?categoryid=&currentRegion=ON&lang=en-CA&page=1&pageSize=24&path=&query=%s&exp=&sortBy=relevance&sortDir=desc", sku)

	req, reqErr := http.NewRequest("GET", link, nil)
	if reqErr != nil {
		return &ProductInfo{}, reqErr
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, resErr := m.Client.Do(req)
	if resErr != nil {
		return &ProductInfo{}, resErr
	}
	defer res.Body.Close()

	body, bodyErr := ioutil.ReadAll(res.Body)
	if bodyErr != nil {
		return &ProductInfo{}, bodyErr
	}

	var resJSON bbresponse.InfoJSON
	unmErr := json.Unmarshal(body, &resJSON)
	if unmErr != nil {
		return &ProductInfo{}, unmErr
	}

	//Print RESJSON Here
	name := resJSON.Products[0].Name
	link = "https://bestbuy.ca" + resJSON.Products[0].ProductURL
	image := resJSON.Products[0].HighResImage
	price := fmt.Sprintf("CAD %.2f", resJSON.Products[0].RegularPrice)
	exec := fmt.Sprintf("%v", time.Since(start))

	return &ProductInfo{name, image, price, sku, link, exec}, nil
}

func (m *Monitor) GetStock(sku string) (*ProductStock, error) {
	start := time.Now()
	m.UpdateClient()

	stock := ProductStock{}
	link := "https://www.bestbuy.ca/ecomm-api/availability/products?accept=application%2Fvnd.bestbuy.standardproduct.v1%2Bjson&accept-language=en-CA&locations=977%7C203%7C931%7C62%7C617%7C927%7C965%7C57%7C938%7C237%7C943%7C932%7C956%7C202%7C200%7C937%7C926%7C795%7C916%7C544%7C910%7C954%7C207%7C233%7C930%7C622%7C223%7C245%7C925%7C985%7C990%7C959%7C949%7C206%7C942&postalCode=M5G2C3&skus=" + sku
	req, reqErr := http.NewRequest("GET", link, nil)
	if reqErr != nil {
		return &stock, reqErr
	}
	for k, v := range headers2 {
		req.Header.Set(k, v)
	}

	res, resErr := m.Client.Do(req)
	if resErr != nil {
		return &stock, resErr
	}
	defer res.Body.Close()

	body, decodeErr := ioutil.ReadAll(res.Body)
	if decodeErr != nil {
		return &stock, decodeErr
	}

	// var stockJSON bbresponse.StockJSON
	// unmErr := json.Unmarshal(body, &stockJSON)
	// if unmErr != nil {
	// 	return &stock, unmErr
	// }

	sVal, _, _, _ := jsonparser.Get(body, "availabilities", "[0]", "shipping", "quantityRemaining")
	stock.Stock, _ = strconv.Atoi(string(sVal))
	sStatus, _, _, _ := jsonparser.Get(body, "availabilities", "[0]", "shipping", "status")
	switch s := string(sStatus); s {
	case "InStock":
		stock.Available = true
	default:
		stock.Available = false
	}
	stock.Exec = fmt.Sprintf("%v", time.Since(start))
	return &stock, nil
}

func (m *Monitor) Initialize() {
	m.ProductInfo = map[string]ProductInfo{}
	for _, sku := range m.Config.SKUs {
		req, reqErr := m.GetInfo(sku)
		if reqErr != nil {
			PrintErr(reqErr)
			m.ProductInfo[sku] = ProductInfo{}
		} else {
			m.ProductInfo[sku] = *req
			GreenMessage(fmt.Sprintf("[ %s ] Initialized %s", *&req.Exec, *&req.Name))
		}
	}
	m.Availability = map[string]bool{}
	for _, sku := range m.Config.SKUs {
		req, reqErr := m.GetStock(sku)
		if reqErr != nil {
			PrintErr(reqErr)
			m.Availability[sku] = false
		} else {
			m.Availability[sku] = req.Available
			GreenMessage(fmt.Sprintf("[ %s ] Initialized Stock for %s", *&req.Exec, sku, *&req.Available))
		}
	}
}

func (m *Monitor) Monitor() {
	i := true
	for i == true {
		for _, sku := range m.Config.SKUs {
			go func(sku string) {
				req, reqErr := m.GetStock(sku)
				if reqErr != nil {
					PrintErr(reqErr)
					return
				}
				YellowMessage(fmt.Sprintf("[ %s ] Monitoring %s", req.Exec, m.ProductInfo[sku].Name))

				if req.Available == true && m.Availability[sku] == false {
					m.Availability[sku] = true
					go m.SendEmbed(sku, *req)
					GreenMessage(fmt.Sprintf("%s Restocked! %v", sku, req.Available))
				} else if req.Available == false && m.Availability[sku] == true {
					m.Availability[sku] = true
					RedMessage(fmt.Sprintf("%s Out of stock! %s", sku, req.Available))
				}
				return
			}(sku)
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func (m *Monitor) SendEmbed(sku string, s ProductStock) {
	prod := m.ProductInfo[sku]
	for _, webhook := range m.Config.Webhooks {
		go func(webhook *godiscord.Webhook, p *ProductInfo, s *ProductStock) {
			e := godiscord.NewEmbed(p.Name, fmt.Sprintf("__**Price:**__ %s", p.Price), p.Link)
			e.AddField("Availability", "In Stock", true)
			e.AddField("Stock", fmt.Sprintf("%d", s.Stock), true)
			e.SetFooter(webhook.Text, webhook.IconURL)
			e.SetColor("#16A085")
			e.SetAuthor("BestBuy CA Monitor", "", "")
			e.SetThumbnail(p.Image)
			e.SendToWebhook(webhook.URL)
		}(&webhook, &prod, &s)
	}
}

func main() {
	m, mErr := NewMonitor("./config.json", "")
	if mErr != nil {
		panic(mErr)
	}

	m.Initialize()
	m.Monitor()
}
