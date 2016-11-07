package main

import (
  "fmt"
  "encoding/json"
  "net/url"
)

type SearchResponse struct {
  Message string `json:"message,omitempty"`
  Product *Product `json:"product,omitempty"`
  Products []Product `json:"products,omitempty"`
  Total int `json:"total,omitempty"`
}

type Field struct {
  AltImages []string `json:"alt_images,omitempty"`
  ExtraInfo string `json:"extra_info,omitempty"`
  Image string `json:"image,omitempty"`
  Price string `json:"price,omitempty"`
  Text string `json:"text,omitempty"`
  Value string `json:"value,omitempty"`
  Weight string  `json:"weight,omitempty"`
  SubFields map[string][]Field `json:"dep,omitempty"`
}

type Product struct {
  URL string `json:"url,omitempty"`
  Title string `json:"title,omitempty"`
  Description string `json:"description,omitempty"`
  Image string`json:"image,omitempty"`
  Price string`json:"price,omitempty"`
  OriginalPrice string`json:"original_price,omitempty"`
  Categories []string`json:"categories,omitempty"`
  SiteCategories []string`json:"site_categories,omitempty"`
  RequiredFieldNames []string`json:"required_field_names,omitempty"`
  RequiredFieldValues map[string][]Field`json:"required_field_values,omitempty"`
  SiteId string`json:"site_id,omitempty"`
  AllCategories []string`json:"all_categories,omitempty"`
  DeepCategories []string`json:"deep_categories,omitempty"`
  MD5 string`json:"md5,omitempty"`
  Status string `json:"status,omitempty"`
  OriginalURL string `json:"original_url,omitempty"`
  CleanURL string `json:"clean_url,omitempty"`
  DiscountedPrice string `json:"discounted_price,omitempty"`
  Returns string `json:"returns,omitempty"`
  PickupSupport bool `json:"pickup_support,omitempty"`
}

type Filter struct {
  Keywords string `json:"keywords,omitempty"`
  SiteIds []string `json:"site_ids,omitempty"`
  Genders []string `json:"genders,omitempty"`
  Categories []string `json:"categories,omitempty"`
  PriceRanges []PriceRange `json:"price_ranges,omitempty"`
}

type PriceRange struct {
  From int `json:"from"`
  To int `json:"to"`
  Currency string `json:"currency"`
}

type SearchQuery struct {
  CurrentPage int `json:"page,omitempty"`
  Per int `json:"per_page,omitempty"`
  Filter Filter `json:"filter,omitempty"`
  Sort string `json:"sort,omitempty"`
  Country string `json:"destination_country,omitempty"`
  client *TwoTap `json:"-"`
}

type GetProductQuery struct {
  MD5 string `json:"product_md5,omitempty"`
  SiteID string `json:"site_id,omitempty"`
  Country string `json:"destination_country,omitempty"`
  client *TwoTap `json:"-"`
}

func (tt *TwoTap) NewSearchQuery() *SearchQuery {
  f := Filter{}
  return &SearchQuery{
    Filter: f,
    Country: tt.DefaultDestination,
    CurrentPage: 1,
    Per: tt.DefaultPageSize,
    client: tt}
}

func (tt *TwoTap) SearchProduct(keywords string) *SearchQuery {
  f := Filter{Keywords: keywords}
  return &SearchQuery{
    Filter: f,
    Country: tt.DefaultDestination,
    CurrentPage: 1,
    Per: tt.DefaultPageSize,
    client: tt}
}

func (q *SearchQuery) FromSites(siteIds []string) *SearchQuery{
  q.Filter.SiteIds = siteIds
  return q
}

func (q *SearchQuery) WithGenders(genders []string) *SearchQuery{
  q.Filter.Genders = genders
  return q
}

func (q *SearchQuery) InRange(priceRanges []PriceRange) *SearchQuery{
  q.Filter.PriceRanges = priceRanges
  return q
}

func (q *SearchQuery) WithCategories(categories []string) *SearchQuery{
  q.Filter.Categories = categories
  return q
}

func (q *SearchQuery) WithFilter(filter Filter) *SearchQuery{
  q.Filter = filter
  return q
}

func (q *SearchQuery) WithCountry(country string) *SearchQuery{
  q.Country = country
  return q
}

func (q *SearchQuery) Page(page int) *SearchQuery{
  q.CurrentPage = page
  return q
}

func (q *SearchQuery) PerPage(per int) *SearchQuery{
  q.Per = per
  return q
}

func (q *SearchQuery) NextPage() *SearchQuery{
  q.CurrentPage += 1
  return q
}

func (q *SearchQuery) Do() *SearchResponse {
  resp, err := q.client.Post("product/search", q)

  if err != nil {
    fmt.Println(err.Error())
    return nil
  }

  searchResponse := SearchResponse{}
  json.Unmarshal(resp, &searchResponse)
  return &searchResponse
}

func (r *SearchResponse) String() (string, error) {
  bytez, err := json.MarshalIndent(r,"","  ")
  return string(bytez), err
}

func (tt *TwoTap) GetProduct(md5 string, siteId string) *GetProductQuery {
  return &GetProductQuery{
    MD5: md5,
    SiteID: siteId,
    Country: tt.DefaultDestination,
    client: tt}
}

func (q *GetProductQuery) WithCountry(country string) *GetProductQuery{
  q.Country = country
  return q
}

func (q *GetProductQuery) Do() *SearchResponse {
  query := fmt.Sprint("product_md5=",q.MD5,"&site_id=",q.SiteID,"&destination_country=",url.QueryEscape(q.Country))

  resp, err := q.client.Get("product",query)
  if err != nil {
    fmt.Println(err.Error())
    return nil
  }

  searchResponse := SearchResponse{}
  json.Unmarshal(resp, &searchResponse)
  return &searchResponse
}

// func (tt *TwoTap) GetProductFull(MD5 string, siteId string, country string) *SearchResponse {
//   if country == "" {
//     country = tt.DefaultDestination
//   }
//   query := fmt.Sprint("product_md5=",MD5,"&site_id=",siteId,"&destination_country=",url.QueryEscape(country))
//
//   resp, err := tt.Get("product",query)
//   if err != nil {
//     fmt.Println(err.Error())
//     return nil
//   }
//
//   searchResponse := SearchResponse{}
//   json.Unmarshal(resp, &searchResponse)
//   return &searchResponse
// }
