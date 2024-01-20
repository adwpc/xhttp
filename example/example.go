package main

import (
	"fmt"

	"github.com/adwpc/xhttp"
	"github.com/buger/jsonparser"
)

type Json struct {
	Args    map[string]string
	Headers map[string]string
	Origin  string
	Url     string
}

func simpleGetRespToString() {
	fmt.Println("--------simpleGetRespToString-----------")
	ip, err := xhttp.New().Get().RespToString("http://httpbin.org/get")

	if err != nil {
		panic(err)
	}
	fmt.Println(ip)
	v, _, _, _ := jsonparser.Get([]byte(ip), "headers", "Host")
	fmt.Println(string(v))
}

func simpleGetRespJsonKey() {
	fmt.Println("--------simpleGetRespJsonKey-----------")
	host, err := xhttp.New().Get().RespGetJsonKey("http://httpbin.org/get", "origin")

	if err != nil {
		panic(err)
	}
	fmt.Println(string(host))
}

func customGetRespToJson() {
	fmt.Println("--------customGetRespToJson-----------")
	var j Json
	err := xhttp.NewWithOption(5000, 5000, 10000, "").Get().AddHeader("a", "b").AddParam("c", "d").SetBody("{\"e\":\"f\"}").RespToJson("http://httpbin.org/get", &j)

	if err != nil {
		panic(err)
	}
	fmt.Println(j)
}

func multiCustomPostRespToJson() {
	fmt.Println("--------multiCustomPostRespToJson-----------")
	headers := make(map[string]string)
	headers["a"] = "b"
	headers["c"] = "d"

	params := make(map[string]string)
	params["q"] = "ip"

	var j Json
	err := xhttp.New().Post().AddHeaders(headers).AddParams(params).SetBody("{\"e\":\"f\"}").RespToJson("http://httpbin.org/post", &j)

	if err != nil {
		panic(err)
	}
	fmt.Println(j)
}
func main() {
	simpleGetRespToString()
	simpleGetRespJsonKey()
	customGetRespToJson()
	multiCustomPostRespToJson()
}
