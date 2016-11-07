package main

type TwoTap struct {
  PublicToken string
  PrivateToken string
  TestMode bool
  DefaultDestination string
  DefaultPageSize int
}

type TTQuery interface {
  Do() TTResponse
}

type TTResponse interface {
  String() (string,error)
}

func NewTwoTap(publicToken string, privateToken string, defaultDestination string, test bool) (*TwoTap, error) {
  if defaultDestination == "" {
    defaultDestination = "United States of America"
  }

  return &TwoTap{publicToken,privateToken,test,defaultDestination,10}, nil
}
