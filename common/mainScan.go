package common

import (
	"sync"
)

func MainWorker(wg *sync.WaitGroup,task chan Task,results []Result)  {
	for tempTask :=range task{
		if tempTask.ServiceName==""{
			close(task)
		}else {
			tempResult,isUnAuth:=scanWorker(tempTask)
			if isUnAuth{
				results=append(results,tempResult)
			}
		}
	}
	wg.Done()
}
func scanWorker(task Task)  (Result,bool){
	if _,ok:=ScanFuncMap[task.ServiceName];ok{
		tempFunc:=ScanFuncMap[task.ServiceName]
		resultStr,isUnAuth:=tempFunc(task.IpPort)
		if isUnAuth{
			return Result{
				IpPort:   task.IpPort,
				UserPass: resultStr,
			},isUnAuth
		}
	}
	return Result{},false
}


