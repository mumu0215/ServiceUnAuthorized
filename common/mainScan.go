package common

import (
	"errors"
	"sync"
)



func MainWorker(wg *sync.WaitGroup,task chan Task,results []Result)  {
	for tempTask :=range task{
		if tempTask.ServiceName==""{
			close(task)
		}else {
			tempResult,err:=scanWorker(tempTask)
			if err==nil{
				results=append(results,tempResult)
			}
		}
	}
	wg.Done()
}
func scanWorker(task Task)  (Result,error){

	return Result{},errors.New("unauthorized failed")
}


