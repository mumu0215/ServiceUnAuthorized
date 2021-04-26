package utill

import (
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"runtime"
	"src/common"
	"strings"
)
var json=jsoniter.ConfigCompatibleWithStandardLibrary
var NumOfCpu=runtime.NumCPU()

func ParseJsonInput(fileName string) (common.ServiceList,error) {
	dataJson,err:=ioutil.ReadFile(fileName)
	if err!=nil{
		return common.ServiceList{},err
	}
	var myServiceList common.ServiceList
	err=json.UnmarshalFromString(string(dataJson),&myServiceList)
	if err!=nil{
		return common.ServiceList{},err
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
	temp:=strings.Split(passwordStr,"\n")
	passwordList= func() []string{
		var temp1 []string
		for _,n :=range temp{
			temp1=append(temp1,strings.TrimSpace(n))
		}
		return temp1
	}()
	return passwordList,nil
}

func GenerateTasks(serviceAll common.ServiceList,isBrute bool) []common.Task {   //读入解析后的json，解析成待扫描的数据格式
	var taskList []common.Task
	for _,ser:=range serviceAll.ServiceList{
		serviceName:=ser.Service
		for _,t:=range ser.IpPortList{
			temp:=common.Task{
				ServiceName: serviceName,
				IpPort:      t,
				IsBrute:     isBrute,
			}
			taskList=append(taskList,temp)
		}
	}
	return taskList
}