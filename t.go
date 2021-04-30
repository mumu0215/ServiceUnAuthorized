package main

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"time"
)
var Client *http.Client
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
func CheckAtlassianCrowd(ipPort string)(string,bool){   //  /crowd/admin/uploadplugin.action  400即存在
	statusCode,_,isSend:=sendGetRequest("http://"+ipPort+"/crowd/admin/uploadplugin.action")
	if isSend{
		if statusCode==400{
			return ipPort+" has Atlassian-Crowd unAuthorized which can lead to rce", true
		}
	}
	return "", false
}
func main(){
	Client=&http.Client{
		Timeout:time.Duration(5)*time.Second,
		Transport: &http.Transport{
			//参数未知影响，目前不使用
			//TLSHandshakeTimeout: time.Duration(timeout) * time.Second,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},}
	fmt.Println(CheckAtlassianCrowd("47.112.108.86:443"))
}
func CheckDubbo(ipPort string)(string,bool){     //存在web和协议端口两种测试方式 Dubbo Admin
	//Authorization: Basic Z3Vlc3Q6Z3Vlc3Q=

	author:=[]string{"Basic Z3Vlc3Q6Z3Vlc3Q=","Basic cm9vdDpyb290"}
	for index,auth:=range author{
		req,err:=http.NewRequest("GET","http://"+ipPort+"/",nil)
		if err!=nil{
			return "", false
		}
		req.Header.Add("Authorization",auth)
		req.Header.Add("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:78.0) Gecko/20100101 Firefox/78.0")
		respond,err:=Client.Do(req)
		if err!=nil{
			return "", false
		}
		defer respond.Body.Close()
		temp,err:=goquery.NewDocumentFromReader(respond.Body)
		if respond.StatusCode==200 && strings.Contains(temp.Text(),"Dubbo Admin"){
			if index==0{
				return ipPort+" has Dubbo unAuthorized, guest/guest", true
			}else {
				return ipPort+" has Dubbo unAuthorized, root/root", true
			}
		}
	}
	return "", false
}