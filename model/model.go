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
type History struct {
	UserOpenId string `json:"userOpenId"`
	From   string `json:"from"`
	To string `json:"to"`
	Word  string `json:"word"`
	Timestamp    string `json:"timestamp";`
	ImageUrl   string `json:"imageurl";`
	Addr   string `json:"addr";`
	EthAddr   string `json:"ethaddr";`
	Height string `json:"height"`

}
type Pay struct {
	Paid bool `json:"paid"`
	Upload_error  bool  `json:"upload_error"`
	From string  `json:"from"`
	To string  `json:"to"`
	Word string  `json:"word"`
	Success bool `json:"success"`
	EthAddr string `json:"ethaddr"`
	Height string `json:"height"`
}
