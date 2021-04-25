package utill

import (
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"strings"
)
var json=jsoniter.ConfigCompatibleWithStandardLibrary

type Service struct {
	Service string `json:"service"`
	IpPortList []string `json:"ip_port_list"`
}
type ServiceList struct {
	ServiceList []Service `json:"service_list"`
}

func ParseJsonInput(fileName string) (ServiceList,error) {
	dataJson,err:=ioutil.ReadFile(fileName)
	if err!=nil{
		return ServiceList{},err
	}
	var myServiceList ServiceList
	err=json.UnmarshalFromString(string(dataJson),&myServiceList)
	if err!=nil{
		return ServiceList{},err
	}
	return myServiceList,nil
}
func ParsePassword(fileName string) ([]string,error) {
	dataPassword,err:=ioutil.ReadFile(fileName)
	if err!=nil{
		return []string{},err
	}
	var passwordList []string
	passwordStr:=strings.TrimSpace(string(dataPassword))
	passwordList=strings.Split(passwordStr,"\n")
	return passwordList,nil
}
