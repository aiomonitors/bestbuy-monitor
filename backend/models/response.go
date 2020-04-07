package models

type Response struct {
	JSONGraph JSONGraph `json:"jsonGraph"`
}

type JSONGraph struct {
	Shop Shop `json:"shop"`
}

type Shop struct {
	Buttonstate Buttonstate `json:"buttonstate"`
}

type Buttonstate struct {
	V5 V5 `json:"v5"`
}

type V5 struct {
	Item Item `json:"item"`
}

type Item struct {
	Skus map[string]Skus `json:"skus"`
}

type Skus struct {
	Conditions Conditions `json:"conditions"`
}

type Conditions struct {
	None None `json:"NONE"`
}

type None struct {
	DestinationZipCode DestinationZipCode `json:"destinationZipCode"`
}

type DestinationZipCode struct {
	The11581 The11581 `json:"11581"`
}

type The11581 struct {
	StoreID StoreID `json:"storeId"`
}

type StoreID struct {
	The950 The950 `json:"950"`
}

type The950 struct {
	Context Context `json:"context"`
}

type Context struct {
	Cyp Cyp `json:"cyp"`
}

type Cyp struct {
	AddAll AddAll `json:"addAll"`
}

type AddAll struct {
	False False `json:"false"`
}

type False struct {
	Type    string `json:"$type"`
	Value   Value  `json:"value"`
	Expires int64  `json:"$expires"`
}

type Value struct {
	ButtonStateResponseInfos []ButtonStateResponseInfo `json:"buttonStateResponseInfos"`
}

type ButtonStateResponseInfo struct {
	SkuID       string `json:"skuId"`
	ButtonState string `json:"buttonState"`
	DisplayText string `json:"displayText"`
}
