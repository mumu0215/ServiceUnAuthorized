package main

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
)
var Client *http.Client
func main(){
	Client=&http.Client{
		Timeout:time.Duration(5)*time.Second,
		Transport: &http.Transport{
			//参数未知影响，目前不使用
			//TLSHandshakeTimeout: time.Duration(timeout) * time.Second,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},}
	fmt.Println(CheckJenkins("39.105.23.173:8080"))
}
func CheckJenkins(ipPort string)(string,bool){   //web port
	statusCode,_,isSend:=sendGetRequest("http://"+ipPort+"/script")
	if !isSend{
		return "", false
	}
	if statusCode==200{
		return ipPort+" has unAuthorized", true
	}
	return "", false
}
func sendGetRequest(url string)(int,string,bool){
	req,err:=http.NewRequest("GET",url,nil)
	if err!=nil{
		return 0,"", false
	}
	respond,err:=Client.Do(req)
	if err!=nil{
		return respond.StatusCode,"", false
	}
	defer respond.Body.Close()
	temp,err:=goquery.NewDocumentFromReader(respond.Body)
	return respond.StatusCode,temp.Text(),true
}
