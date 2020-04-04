# Bestbuy Monitor
Given a list of links to monitor, this program monitors the availability of products on [Bestbuy.com](https://www.bestbuy.com).
The program sends a notification to a discord webhook when a product changes from OOS to In stock. 
Built because of high demand to purchase a Nintendo Switch during the COVID-19 outbreak.

## Requirements
Go 1.1+

## Instructions
1. Clone the repo using `git clone https://github.com/aiomonitors/bestbuy-monitor`
2. Setup your `config.json` file. It should look like the file below. Links should be an array of links, not a string. 
```json
{
    "webhooks" : [
        {
            "webhook" : "WEBHOOK URL HERE",
            "icon_url" : "Link to icon here, or leave as empty string",
            "text" : "Name to display in footer, or leave empty"
        }
    ],
    "color" : "#F1B379",
    "links" : ["https://www.bestbuy.com/site/nintendo-switch-32gb-console-neon-red-neon-blue-joy-con/6364255.p?skuId=6364255"]
}
```
3. To start the program, use `go run main.go --config=NAME OF CONFIG FILE --proxy=NAME OF PROXY FILE`. The proxy file should be a `.txt` file 
containing proxies in the format `https://ip:port:username:password`. The `go.mod` file should automatically download all required dependencies, if not run 
`go mod download`
4. Happy monitoring!


If you have any questions, feel free to contact me on: [Twitter](https://twitter.com/aiomonitors) or shoot me an email at [navr@discoders.us](mailto:navr@discoders.us)
My personal site: [shihab.dev](https://www.shihab.dev)