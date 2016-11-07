package main

import (
  "net/http"
  "encoding/json"
  "bytes"
  "errors"
  "io/ioutil"
  "fmt"
)

const (
  BASE_URL = "https://api.twotap.com/v1.0/"
  JSON_CONTENT = "application/json"
)

var client *http.Client

func getClient() *http.Client{
  if client == nil{
    client = &http.Client{}
  }
  return client
}

func (tt *TwoTap) Post(path string, data interface{}) ([]byte,error) {
  bytez, err := json.Marshal(data)
  debugBytez, err := json.MarshalIndent(data,"","  ")

  fmt.Println(string(debugBytez))

  if err != nil {
    return []byte{}, err
  }

  resp,err := getClient().Post(BASE_URL + tt.appendToken(path,""), JSON_CONTENT, bytes.NewReader(bytez))
  if err != nil {
    return []byte{}, err
  }

  err = catchErrors(resp)
  if err != nil {
    return []byte{}, err
  }

  content,err := ioutil.ReadAll(resp.Body)
  defer resp.Body.Close()

  return content, nil
}

func (tt *TwoTap) Get(path string, query string) ([]byte,error) {
  fmt.Println(BASE_URL + tt.appendToken(path, query))
  resp,err := getClient().Get(BASE_URL + tt.appendToken(path, query))
  if err != nil {
    return []byte{}, err
  }

  err = catchErrors(resp)
  if err != nil {
    return []byte{}, err
  }

  content,err := ioutil.ReadAll(resp.Body)
  defer resp.Body.Close()

  return content, nil
}

func catchErrors(resp *http.Response) error {
  if resp.StatusCode > 399 {
    body,_ := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()

    return errors.New(resp.Status + " : " + string(body))
  }
  return nil
}

func (tt *TwoTap) appendToken(path, query string) string {
  if path == "purchase/confirm" ||
      path == "purchase/history" ||
      path == "wallet/user_token" ||
      path == "coupons" {
    if query == "" {
      query = "private_token=" + tt.PrivateToken
    }else {
      query += "&private_token=" + tt.PrivateToken
    }
  } else {
    if query == "" {
      query = "public_token=" + tt.PublicToken
    }else {
      query += "&public_token=" + tt.PublicToken
    }
  }
  return path + "?" + query
}
