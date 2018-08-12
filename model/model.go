package model

type Node struct {
	UserOpenId string `json:"userOpenId"`
	From   string `json:"form"`
	To string `json:"to"`
	Word  string `json:"word"`
	Timestamp    string `json:"timestamp";`
	ImageUrl   string `json:"imageurl";`
	Addr   string `json:"addr";`
}
