// Code generated by goctl. DO NOT EDIT.
package types

type ArticleReq struct {
	Title  string `json:"title" redis:"title"`
	Link   string `json:"link" redis:"link"`
	Poster string `json:"poster" redis:"poster"`
	Time   string `json:"time" redis:"time"`
	Votes  int    `json:"votes" redis:"votes"`
}

type ArticleResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
