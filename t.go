package main

import (
	"fmt"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func main(){
	//a:=common.CheckMemCache("10.4.24.166:11211")
	//fmt.Println(a)
	res,_:=http.Get("http://118.31.34.69:9200/_cat")
	b,_:=goquery.NewDocumentFromReader(res.Body)
	fmt.Println(dealWithEla(b.Text()))
}
func dealWithEla(s string)bool{
	temp:=strings.Split(s,"\n")[1:]
	fmt.Println(temp)
	for _,i:=range temp{
		if i=="/_cat/shards"{
			fmt.Println(i)
			return true
		}
	}
	return false
}