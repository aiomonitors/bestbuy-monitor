package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aiomonitors/bestbuy-monitor/backend/models"
	"github.com/aiomonitors/godiscord"
	proxymanager "github.com/aiomonitors/goproxymanager"
	"github.com/fatih/color"
)

var headers = map[string]string{
	"authority":       "www.bestbuy.com",
	"cache-control":   "max-age=0",
	"user-agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
	"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"accept-language": "en-US,en;q=0.9",
}

type Config struct {
	Webhooks []godiscord.Webhook `json:"webhooks"`
	Color    string              `json:"color"`
	SKUs     []string            `json:"skus"`
}

type ProductInfo struct {
	Name         string `json:"name"`
	SKU          string `json:"sku"`
	Availability bool   `json:"availability"`
	Link         string `json:"link"`
	Image        string `json:"image"`
	Exec         string `json:"exec"`
}

type Monitor struct {
	Config       Config                    `json:"config"`
	Availability map[string]bool           `json:"availability`
	ProductInfo  map[string]ProductInfo    `json:"product_info"`
	Manager      proxymanager.ProxyManager `json:"manager"`
	UseProxies   bool                      `json:"use_proxies"`
	Client       http.Client               `json:"client"`
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

func (m *Monitor) GetInfo(sku string) (*ProductInfo, error) {
	m.UpdateClient()
	start := time.Now()
	//Sets up http client with proxies if monitor uses proxies
	link := fmt.Sprintf("https://www.bestbuy.com/site/%s.p", sku)
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
	sku = document.Find("div.sku.product-data > span.product-data-value.body-copy").Text()
	buttonText := document.Find(".fulfillment-add-to-cart-button > div").Find("button").Text()
	image, _ := document.Find("div.primary-image-container > button > img").Attr("src")

	elapsed := fmt.Sprintf("%s", time.Since(start))

	if buttonText != "" && buttonText == "Add to Cart" {
		return &ProductInfo{name, sku, false, link, image, elapsed}, nil
	}
	return &ProductInfo{name, sku, false, link, image, elapsed}, nil
}

func (m *Monitor) GetStockJSON() (map[string]bool, error) {
	m.UpdateClient()
	defRes := map[string]bool{}
	rand.Seed(time.Now().UnixNano())
	//[["shop","buttonstate","v5","item","skus",[6364255,6364253,6401728,6342914,6342915,6289641,5761912,2588445,4503702,6084400],"conditions","NONE","destinationZipCode",11581,"storeId",950,"context","cyp","addAll","false"]]
	link := "https://www.bestbuy.com/api/tcfb/model.json?paths=%5B%5B%22shop%22,%22buttonstate%22,%22v5%22,%22item%22,%22skus%22,%5B" + strings.Join(m.Config.SKUs, ",") + "%5D,%22conditions%22,%22NONE%22,%22destinationZipCode%22,11581,%22storeId%22,950,%22context%22,%22cyp%22,%22addAll%22,%22false%22%5D%5D&method=post" + fmt.Sprintf("&rand=%s", rand.Int63n(12312371928371283))

	req, reqErr := http.NewRequest("GET", link, nil)
	if reqErr != nil {
		return defRes, reqErr
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, resErr := m.Client.Do(req)
	if resErr != nil {
		return defRes, resErr
	}
	defer res.Body.Close()

	var resJSON models.Response
	body, bodyReadErr := ioutil.ReadAll(res.Body)
	if bodyReadErr != nil {
		return defRes, bodyReadErr
	}

	unmErr := json.Unmarshal(body, &resJSON)
	if unmErr != nil {
		return defRes, unmErr
	}

	availability := map[string]bool{}
	for _, sku := range m.Config.SKUs {
		buttonStatus := resJSON.JSONGraph.Shop.Buttonstate.V5.Item.Skus[sku].Conditions.None.DestinationZipCode.The11581.StoreID.The950.Context.Cyp.AddAll.False.Value.ButtonStateResponseInfos[0].ButtonState
		if buttonStatus == "SOLD_OUT" {
			availability[sku] = false
		} else {
			availability[sku] = true
		}
	}

	return availability, nil
}

func (m *Monitor) InitializeProducts() {
	m.ProductInfo = map[string]ProductInfo{}
	for _, sku := range m.Config.SKUs {
		req, reqErr := m.GetInfo(sku)
		if reqErr != nil {
			PrintErr(reqErr)
		}
		m.ProductInfo[sku] = *req
		GreenMessage(fmt.Sprintf("[ %s ] Initialized info for %s", req.Exec, req.Name))
	}
}

func (m *Monitor) Initialize() {
	start := time.Now()
	m.InitializeProducts()
	sj, _ := m.GetStockJSON()
	m.Availability = sj
	GreenMessage(fmt.Sprintf("[ %s ] Initialized", time.Since(start)))
}

func (m *Monitor) Monitor() {
	i := true
	for i == true {
		time.Sleep(time.Millisecond * 1000)

		start := time.Now()
		req, reqErr := m.GetStockJSON()
		if reqErr != nil {
			PrintErr(reqErr)
			continue
		}

		for sku, avail := range req {
			go func(sku string, avail bool) {
				if avail == true && m.Availability[sku] == false {
					m.Availability[sku] = true
					go m.SendEmbed(sku)
					GreenMessage(fmt.Sprintf("%s restocked!\n", sku))
				} else if avail == false && m.Availability[sku] == true {
					m.Availability[sku] = false
					RedMessage(fmt.Sprintf("%s is out of stock", sku))
				}
			}(sku, avail)
		}
		color.Blue("[ %s ] Monitored in [ %v ]", GetTime(), time.Since(start))
	}
}

func (m *Monitor) SendEmbed(sku string) {
	for _, webhook := range m.Config.Webhooks {
		go func(sku string, webhook godiscord.Webhook) {
			p := m.ProductInfo[sku]
			e := godiscord.NewEmbed(p.Name, fmt.Sprintf("In Stock: True"), p.Link)
			e.AddField("ATC", fmt.Sprintf("[ATC](https://api.bestbuy.com/click/-/%s/cart)", sku), false)
			e.SetFooter(webhook.Text, webhook.IconURL)
			e.SetColor("#16A085")
			e.SetAuthor("Best Buy Monitor [ Backend ]", "", "")
			e.SetThumbnail(p.Image)
			e.SendToWebhook(webhook.URL)
		}(sku, webhook)
	}
}

func main() {
	m, _ := NewMonitor("./config.json", "")
	m.Initialize()
	m.Monitor()
}
