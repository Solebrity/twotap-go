package main

import (
  "fmt"
  "encoding/json"
  "net/url"
)

type CartQuery struct {
  Products []string `json:"products,omitempty"`
  FinishedURL string `json:"finished_url,omitempty"`
  FinishedProductAttributesFormat string `json:"finished_products_attributes_format,omitempty"`
  Notes map[string]interface{} `json:"notes,omitempty"`
  TestMode string `json:"test_mode,omitempty"`
  CacheTime int `json:"cache_time,omitempty"`
  Country string `json:"destination_country,omitempty"`
  client *TwoTap
}

type CartStatusQuery struct {
  CartID string `json:"cart_id"`
  ProductAttributesFormat string `json:"products_attributes_format,omitempty"`
  TestMode string `json:"test_mode,omitempty"`
  Country string `json:"destination_country,omitempty"`
  client *TwoTap `json:"-"`
}

type CartResponse struct {
  ID string `json:"cart_id,omitempty"`
  Message string `json:"message,omitempty"`
  Description string `json:"description,omitempty"`
  TestMode string `json:"-"`
  Country string `json:"-"`
  client *TwoTap `json:"-"`
}

type CartStatusResponse struct {
  Message string `json:"message,omitempty"`
  Description string `json:"description,omitempty"`
  UnknownURLS string `json:"unknown_urls,omitempty"`
  Notes string `json:"notes,omitempty"`
  Country string `json:"country,omitempty"`
  DestinationCountry string `json:"destination_country,omitempty"`
  Sites map[string]Site `json:"sites,omitempty"`
}

type Site struct {
  Info SiteInfo `json:"info,omitempty"`
  Coupons bool `json:"coupon_support,omitempty"`
  GiftCards bool `json:"gift_card_support,omitempty"`
  CheckoutTypes []string `json:"checkout_support,omitempty"`
  ShippingCountries []string `json:"shipping_countries_support,omitempty"`
  BillingCountries []string `json:"billing_countries_support,omitempty"`
  ShipingOptions map[string]string `json:"shipping_options,omitempty"`
  Cart map[string]Product `json:"add_to_cart,omitempty"`
  Failures map[string]Product `json:"failed_to_add_to_cart,omitempty"`
  Returns string `json:"returns,omitempty"`
  RemovedProducts []string `json:"removed_products,omitempty"`
  Prices CheckoutPrices `json:"prices,omitempty"`
  Details CheckoutDetails `json:"details,omitempty"`
  StatusMessages []string `json:"status_messages,omitempty"`
  StatusReason string `json:"status_reason,omitempty"`
  RemoteState map[string]string `json:"remote_stat,omitempty"`
  OrderId string `json:"order_id,omitempty"`
  Products map[string]Product `json:"products,omitempty"`
}

type SiteInfo struct {
  Logo string `json:"logo,omitempty"`
  Name string `json:"name,omitempty"`
  URL string `json:"url,omitempty"`
}


func (tt *TwoTap) NewCart(products []string) *CartQuery {
  var test string
  if tt.TestMode {
    test = "fake_confirm"
  }
  return &CartQuery{TestMode: test,
    Products:products,
    CacheTime: 300,
    Country: tt.DefaultDestination,
    client: tt}
}

func (tt *TwoTap) GetCartStatus(cartID string) *CartStatusQuery {
  var test string
  if tt.TestMode {
    test = "fake_confirm"
  }
  return &CartStatusQuery{TestMode: test,
    Country: tt.DefaultDestination,
    CartID: cartID,
    client: tt}
}

func (r *CartResponse) Status() *CartStatusResponse {
  q := &CartStatusQuery{TestMode: r.TestMode,
    Country: r.Country,
    CartID: r.ID,
    client: r.client}

  return q.Do()
}

func (q *CartQuery) AddProduct(product string) *CartQuery {
  q.Products = append(q.Products, product)
  return q
}

func (q *CartQuery) AddProducts(products []string) *CartQuery {
  q.Products = products
  return q
}

func (q *CartQuery) WithFinishedURL(url string) *CartQuery {
  q.FinishedURL = url
  return q
}

func (q *CartQuery) WithFinishedFormat(format string) *CartQuery {
  q.FinishedProductAttributesFormat = format
  return q
}

func (q *CartQuery) WithNotes(notes map[string]interface{}) *CartQuery {
  q.Notes = notes
  return q
}

func (q *CartQuery) WithCacheTime(cache int) *CartQuery {
  q.CacheTime = cache
  return q
}

func (q *CartQuery) WithCountry(country string) *CartQuery {
  q.Country = country
  return q
}

func (q *CartQuery) Do() *CartResponse {
  resp, err := q.client.Post("cart", q)

  if err != nil {
    fmt.Println(err.Error())
    return nil
  }

  cartResponse := CartResponse{}
  json.Unmarshal(resp, &cartResponse)
  cartResponse.client = q.client
  cartResponse.Country = q.Country
  cartResponse.TestMode = q.TestMode
  return &cartResponse
}

func (r *CartResponse) String() (string, error) {
  bytez, err := json.MarshalIndent(r,"","  ")
  return string(bytez), err
}

func (q *CartStatusQuery) WithFormat(format string) *CartStatusQuery {
  q.ProductAttributesFormat = format
  return q
}

func (q *CartStatusQuery) WithCountry(country string) *CartStatusQuery {
  q.Country = country
  return q
}

func (q *CartStatusQuery) Do() *CartStatusResponse {
  query := fmt.Sprint("cart_id=",q.CartID,"&destination_country=",url.QueryEscape(q.Country))

  if q.ProductAttributesFormat != "" {
    query += "&products_attributes_format=" + q.ProductAttributesFormat
  }

  if q.TestMode != "" {
    query += "&test_mode=" + q.TestMode
  }

  resp, err := q.client.Get("cart/status",query)
  if err != nil {
    fmt.Println(err.Error())
    return nil
  }

  cartStatusResponse := CartStatusResponse{}
  json.Unmarshal(resp, &cartStatusResponse)
  return &cartStatusResponse
}

func (r *CartStatusResponse) String() (string, error) {
  bytez, err := json.MarshalIndent(r,"","  ")
  return string(bytez), err
}
