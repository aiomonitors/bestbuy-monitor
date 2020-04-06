package bbresponse

type StockJSON struct {
	Availabilities []Availability `json:"availabilities"`
}

type Availability struct {
	Pickup                 Pickup   `json:"pickup"`                
	Shipping               Shipping `json:"shipping"`              
	Sku                    string   `json:"sku"`                   
	SellerID               string   `json:"sellerId"`              
	SaleChannelExclusivity string   `json:"saleChannelExclusivity"`
	ScheduledDelivery      bool     `json:"scheduledDelivery"`     
	IsGiftCard             bool     `json:"isGiftCard"`            
	IsService              bool     `json:"isService"`             
}

type Pickup struct {
	Status      string        `json:"status"`     
	Purchasable bool          `json:"purchasable"`
	Locations   []interface{} `json:"locations"`  
}

type Shipping struct {
	Status                 string        `json:"status"`                
	QuantityRemaining      int64         `json:"quantityRemaining"`     
	Purchasable            bool          `json:"purchasable"`           
	LevelsOfServices       []interface{} `json:"levelsOfServices"`      
	OrderLimit             int64         `json:"orderLimit"`            
	RestrictedZoneRegions  []interface{} `json:"restrictedZoneRegions"` 
	HasActiveCountdown     bool          `json:"hasActiveCountdown"`    
	CountdownIsZone        bool          `json:"countdownIsZone"`       
	PreorderInfo           interface{}   `json:"preorderInfo"`          
	IsFreeShippingEligible bool          `json:"isFreeShippingEligible"`
	IsBackorderable        bool          `json:"isBackorderable"`       
}


//Product
// Generated by https://quicktype.io
//
// To change quicktype's target language, run command:
//
//   "Set quicktype target language"

type InfoJSON struct {
	Brand             string        `json:"Brand"`            
	CurrentPage       int64         `json:"currentPage"`      
	Total             int64         `json:"total"`            
	TotalPages        int64         `json:"totalPages"`       
	PageSize          int64         `json:"pageSize"`         
	Products          []Product     `json:"products"`         
	Paths             []Path        `json:"paths"`            
	Facets            []Facet       `json:"facets"`           
	LastSearchDate    string        `json:"LastSearchDate"`   
	RelatedQueries    interface{}   `json:"relatedQueries"`   
	Sscs              []interface{} `json:"sscs"`             
	RelatedCategories interface{}   `json:"relatedCategories"`
	SelectedFacets    []interface{} `json:"selectedFacets"`   
	Resources         interface{}   `json:"resources"`        
	RedirectURL       interface{}   `json:"redirectUrl"`      
	Promotions        interface{}   `json:"promotions"`       
	HasBrandStore     bool          `json:"hasBrandStore"`    
	ProductStatusCode string        `json:"productStatusCode"`
	SscStatusCode     interface{}   `json:"sscStatusCode"`    
	Breadcrumb        []interface{} `json:"breadcrumb"`       
}

type Facet struct {
	Name          string   `json:"name"`         
	SystemName    string   `json:"systemName"`   
	Count         int64    `json:"count"`        
	Filters       []Filter `json:"filters"`      
	IsMultiSelect bool     `json:"isMultiSelect"`
	Order         int64    `json:"order"`        
}

type Filter struct {
	Name       string `json:"name"`      
	Path       string `json:"path"`      
	Count      int64  `json:"count"`     
	IsSelected bool   `json:"isSelected"`
}

type Path struct {
	Name         string `json:"name"`        
	Value        string `json:"value"`       
	SelectPath   string `json:"selectPath"`  
	UnselectPath string `json:"unselectPath"`
	Type         string `json:"type"`        
}

type Product struct {
	AltLangSEOText          interface{} `json:"altLangSeoText"`         
	CategoryName            string      `json:"categoryName"`           
	ContentAggregator       interface{} `json:"contentAggregator"`      
	HighResImage            string      `json:"highResImage"`           
	IsMarketplace           bool        `json:"isMarketplace"`          
	OfferID                 interface{} `json:"offerId"`                
	PrimaryParentCategoryID interface{} `json:"primaryParentCategoryId"`
	Seller                  interface{} `json:"seller"`                 
	SellerID                interface{} `json:"sellerId"`               
	SEOText                 string      `json:"seoText"`                
	Sku                     string      `json:"sku"`                    
	Name                    string      `json:"name"`                   
	RegularPrice            float64     `json:"regularPrice"`           
	SalePrice               float64     `json:"salePrice"`              
	ShortDescription        string      `json:"shortDescription"`       
	ProductType             interface{} `json:"productType"`            
	ThumbnailImage          string      `json:"thumbnailImage"`         
	CustomerRating          float64     `json:"customerRating"`         
	CustomerRatingCount     int64       `json:"customerRatingCount"`    
	CustomerReviewCount     int64       `json:"customerReviewCount"`    
	ProductURL              string      `json:"productUrl"`             
	IsAdvertised            bool        `json:"isAdvertised"`           
	IsClearance             bool        `json:"isClearance"`            
	IsInStoreOnly           bool        `json:"isInStoreOnly"`          
	IsOnlineOnly            bool        `json:"isOnlineOnly"`           
	IsPreorderable          bool        `json:"isPreorderable"`         
	IsVisible               bool        `json:"isVisible"`              
	HasPromotion            bool        `json:"hasPromotion"`           
	IsFrenchCompliant       bool        `json:"isFrenchCompliant"`      
	Ehf                     int64       `json:"ehf"`                    
	CurrentRegion           interface{} `json:"currentRegion"`          
	HideSavings             bool        `json:"hideSavings"`            
	IsPriceEndsLabel        bool        `json:"isPriceEndsLabel"`       
	HideSaleEndDate         bool        `json:"hideSaleEndDate"`        
	SaleEndDate             interface{} `json:"saleEndDate"`            
	PriceUnit               string      `json:"priceUnit"`              
	RequiresAgeVerification bool        `json:"requiresAgeVerification"`
	HasFrenchContent        bool        `json:"hasFrenchContent"`       
}
