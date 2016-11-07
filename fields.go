package main

import (
  "fmt"
  "encoding/json"
)

type ValidationRequest struct {
  CartID string `json:"cart_id,omitempty"`
  Fields map[string]interface{} `json:"flat_fields_input,omitempty"`
  client *TwoTap `json:"-"`
}

type ValidationResponse struct {
  Message string `json:"message,omitempty"`
  Description string `json:"description,omitempty"`
}

func (tt *TwoTap) NewValidationRequest(cartId string, fields map[string]interface{}) *ValidationRequest {
  return &ValidationRequest{CartID: cartId, Fields:fields, client: tt}
}

func (q *ValidationRequest) SetFields(fields map[string]interface{}) *ValidationRequest {
  q.Fields = fields
  return q
}

func (q *ValidationRequest) Do() *ValidationResponse {
  resp, err := q.client.Post("fields_input_validate", q)

  if err != nil {
    fmt.Println(err.Error())
    return nil
  }

  validationResponse := ValidationResponse{}
  json.Unmarshal(resp, &validationResponse)
  return &validationResponse
}

func (r *ValidationResponse) String() (string,error){
  bytez, err := json.MarshalIndent(r,"","  ")
  return string(bytez), err
}
