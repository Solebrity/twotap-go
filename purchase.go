package main

import (
  "fmt"
  "encoding/json"
)

type PurchaseRequest struct {
  CartId string `json:"cart_id,omitempty"`
  CheckoutFields map[string]Checkout `json:"fields_input,omitempty"`
  AffiliateLinks map[string]string `json:"affiliate_links,omitempty"`
  Confirm *PurchaseCallback `json:"confirm,omitempty"`
  Products []string `json:"products,omitempty"`
  Notes map[string]interface{} `json:"notes,omitempty"`
  TestMode string `json:"test_mode,omitempty"`
  Locale string `json:"locale,omitempty"`
  client *TwoTap `json:"-"`
}

type PurchaseCallback struct {
  Method string `json:"method,omitempty"`
  ConfirmURL string `json:"http_confirm_url,omitempty"`
  UpdateURL string `json:"http_update_url,omitempty"`
  SkipConfirm string `json:"skip_confirm,omitempty"`
}

type PurchaseResponse struct {
  ID string `json:"purchase_id,omitempty"`
  Message string `json:"message,omitempty"`
  Description string `json:"description,omitempty"`
  client *TwoTap `json:"-"`
  TestMode string `json:"-"`
}

type Checkout interface {}

type AuthCheckout struct {
  Fields *CheckoutOptions `json:"authCheckout,omitempty"`
  Login *Credentials `json:"login,omitempty"`
  Cart map[string]map[string]interface{} `json:"addToCart,omitempty"`
  Coupons []string  `json:"coupons,omitempty"`
  GiftCard *GiftCard `json:"gift_card,omitempty"`
  ShippingOption string  `json:"shipping_option,omitempty"`
}

type NoAuthCheckout struct {
  Fields *CheckoutOptions `json:"noauthCheckout,omitempty"`
  Cart map[string]map[string]interface{} `json:"addToCart,omitempty"`
  Coupons []string  `json:"coupons,omitempty"`
  GiftCard *GiftCard `json:"gift_card,omitempty"`
  ShippingOption string  `json:"shipping_option,omitempty"`
}

type LocalCheckout struct {
  Fields *CheckoutOptions `json:"localCheckout,omitempty"`
  Cart map[string]map[string]interface{} `json:"addToCart,omitempty"`
  Coupons []string  `json:"coupons,omitempty"`
  GiftCard *GiftCard `json:"gift_card,omitempty"`
  Pickup map[string]string `json:"pickup,omitempty"`
}

type GiftCard struct {
  Number string `json:"number,omitempty"`
  PIN string `json:"pin,omitempty"`
}

type Credentials struct {
  Username string `json:"username,omitempty"`
  Password string `json:"password,omitempty"`
}

type CheckoutOptions struct {
  Email string `json:"email,omitempty"`
  ShippingFirstName string `json:"shipping_first_name,omitempty"`
  ShippingLastName string `json:"shipping_last_name,omitempty"`
  ShippingAddress string `json:"shipping_address,omitempty"`
  ShippingCity string `json:"shipping_city,omitempty"`
  ShippingState string `json:"shipping_state,omitempty"`
  ShippingCountry string `json:"shipping_country,omitempty"`
  ShippingZip string `json:"shipping_zip,omitempty"`
  ShippingPhone string `json:"shipping_telephone,omitempty"`
  BillingFirstName string `json:"billing_first_name,omitempty"`
  BillingLastName string `json:"billing_last_name,omitempty"`
  BillingAddress string `json:"billing_address,omitempty"`
  BillingCity string `json:"billing_city,omitempty"`
  BillingState string `json:"billing_state,omitempty"`
  BillingCountry string `json:"billing_country,omitempty"`
  BillingZip string `json:"billing_zip,omitempty"`
  BillingPhone string `json:"billing_telephone,omitempty"`
  CardType string `json:"card_type,omitempty"`
  CardNumber string `json:"card_number,omitempty"`
  CardName string `json:"card_name,omitempty"`
  CardExpiryYear string `json:"expiry_date_year,omitempty"`
  CardExpiryMonth string `json:"expiry_date_month,omitempty"`
  CardCVV string `json:"cvv,omitempty"`
}

type PurchaseStatusRequest struct {
  PurchaseID string `json:"purchase_id,omitempty"`
  TestMode string `json:"test_mode,omitempty"`
  client *TwoTap `json:"-"`
}

type PurchaseStatusResponse struct {
  PurchaseID string `json:"purchase_id,omitempty"`
  Message string `json:"message,omitempty"`
  Description string `json:"description,omitempty"`
  State string `json:"state,omitempty"`
  Pending bool `json:"pending_confirm,omitempty"`
  CreatedAt string `json:"created_at,omitempty"`
  TotalPrices *CheckoutPrices `json:"total_prices,omitempty"`
  Destination string `json:"destination,omitempty"`
  Notes string `json:"notes,omitempty"`
  Sites map[string]Site `json:"sites,omitempty"`
  UserId string `json:"user_id,omitempty"`
  UsedProfiles map[string]interface{} `json:"used_profiles,omitempty"`
  Expiration int64 `json:"session_finishes_at,omitempty"`
}

type CheckoutDetails struct {
  ShippingEstimate string `json:"shipping_estimate,omitempty"`
  ActivePaymentMethod string `json:"active_payment_method,omitempty"`
  ActiveShippingAddress string `json:"active_shipping_address,omitempty"`
  Coupons map[string]CouponStatus  `json:"coupons,omitempty"`
}

type CheckoutPrices struct{
  FinalPrice string `json:"final_price,omitempty"`
  ShippingPrice string `json:"shipping_price,omitempty"`
  SalesTax string `json:"sales_tax,omitempty"`
  GiftCard string `json:"gift_card_value,omitempty"`
  Coupon string `json:"coupon_value,omitempty"`
}

type CouponStatus struct {
  Status string `json:"status,omitempty"`
}

type PurchaseConfirmationRequest struct {
  PurchaseID string `json:"purchase_id,omitempty"`
  TestMode string `json:"test_mode,omitempty"`
}

type PurchaseConfirmationResponse struct {
  PurchaseID string `json:"purchase_id,omitempty"`
  Message string `json:"message,omitempty"`
  Description string `json:"description,omitempty"`
}

func (tt *TwoTap) NewPurchaseRequest(cartId string, checkoutFields map[string]Checkout) *PurchaseRequest{
  var test string
  if tt.TestMode {
    test = "fake_confirm"
  }
  return &PurchaseRequest{CartId: cartId, CheckoutFields: checkoutFields, client: tt,TestMode: test}
}

func (q *PurchaseRequest) WithAffiliateLinks(links map[string]string) *PurchaseRequest{
  q.AffiliateLinks = links
  return q
}

func (q *PurchaseRequest) WithCallback(cb PurchaseCallback) *PurchaseRequest{
  q.Confirm = &cb
  return q
}

func (q *PurchaseRequest) WithProducts(products []string) *PurchaseRequest{
  q.Products = products
  return q
}

func (q *PurchaseRequest) WithNotes(notes map[string]interface{}) *PurchaseRequest{
  q.Notes = notes
  return q
}

func (q *PurchaseRequest) InLocale(locale string) *PurchaseRequest{
  q.Locale = locale
  return q
}

func (q *PurchaseRequest) Do() *PurchaseResponse {
  resp, err := q.client.Post("purchase", q)

  if err != nil {
    fmt.Println(err.Error())
    return nil
  }

  purchaseResponse := PurchaseResponse{TestMode: q.TestMode, client: q.client}
  json.Unmarshal(resp, &purchaseResponse)
  return &purchaseResponse
}

func (r *PurchaseResponse) String() (string, error) {
  bytez, err := json.MarshalIndent(r,"","  ")
  return string(bytez), err
}

func (tt *TwoTap) NewPurchaseStatusRequest(purchase_id string) *PurchaseStatusRequest {
  var test string
  if tt.TestMode {
    test = "fake_confirm"
  }
  return &PurchaseStatusRequest{PurchaseID:purchase_id,TestMode: test}
}

func (r *PurchaseResponse) Status() *PurchaseStatusResponse{
  req := &PurchaseStatusRequest{TestMode: r.TestMode, client: r.client, PurchaseID: r.ID}
  return req.Do()
}

func (q *PurchaseStatusRequest) Do() *PurchaseStatusResponse{
  query := fmt.Sprint("purchase_id=",q.PurchaseID,"&test_mode=",q.TestMode)

  resp, err := q.client.Get("purchase/status",query)
  if err != nil {
    fmt.Println(err.Error())
    return nil
  }

  purchaseStatusResponse := PurchaseStatusResponse{}
  json.Unmarshal(resp, &purchaseStatusResponse)
  return &purchaseStatusResponse
}

func (r *PurchaseStatusResponse) String() (string, error) {
  bytez, err := json.MarshalIndent(r,"","  ")
  return string(bytez), err
}

func (tt *TwoTap) NewPurchaseConfirmRequest(string purchaseID) *PurchaseConfirmationRequest {
  test := ""
  if tt.TestMode {
    test = "fake_confirm"
  }

  return &PurchaseConfirmationRequest{PurchaseID: purchaseId, TestMode: test}
}

func (r *PurchaseResponse) Confirm() *PurchaseConfirmationResponse {
  req := &PurchaseConfirmationRequest{PurchaseID: r.PurchaseID, TestMode: r.TestMode}
  return req.Do()
}

func (q *PurchaseConfirmationRequest) Do() *PurchaseConfirmationResponse {
  query := fmt.Sprint("purchase_id=",q.PurchaseID,"&test_mode=",q.TestMode)


  resp, err := q.client.Get("purchase/confirm",query)
  if err != nil {
    fmt.Println(err.Error())
    return nil
  }

  purchaseConfirmationResponse := PurchaseConfirmationResponse{}
  json.Unmarshal(resp, &purchaseConfirmationResponse)
  return &purchaseConfirmationResponse
}

func (r *PurchaseConfirmationResponse) String() (string, error) {
  bytez, err := json.MarshalIndent(r,"","  ")
  return string(bytez), err
}
