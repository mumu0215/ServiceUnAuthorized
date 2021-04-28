package utill

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"src/common"
	"sync"
)

var(
	InputFile string
	Thread int
	PassWordFile string
	IsBrute bool
	PassWordList []string
)
func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func FlagMain(c *cli.Context) error{
	if c.IsSet("input"){  //校验文件存在
		if !Exists(InputFile){
			return errors.New("input file not exist")
		}
	}
	var myService common.ServiceList
	var err error
	myService,err=ParseJsonInput(InputFile)
	if err!=nil{
		return err
	}
	fmt.Println(myService.ServiceList)
	if IsBrute{     //启用爆破模块，主要校验线程和密码文件名合法性
		if c.IsSet("thread"){   //检验线程数合法
			if Thread<=0{
				return errors.New("illegal thread number")
			}
		}
		if c.IsSet("password"){  //校验文件存在
			if !Exists(PassWordFile){
				return errors.New("password file not exist")
			}
			PassWordList,err=ParsePassword(PassWordFile)
			if err!=nil{
				return err
			}
			fmt.Println(PassWordList)
		}
	}
	//处理任务采用生产/消费模式，使用channel进行数据交互
	var myTaskList []common.Task
	taskChannel:=make(chan common.Task)
	var results []common.Result
	myTaskList=GenerateTasks(myService,IsBrute)
	wg:=&sync.WaitGroup{}
	//numOfTask:=len(myTaskList)

	//启用协程接受任务
	for i:=0;i<Thread;i++{
		wg.Add(1)
		go common.MainWorker(wg,taskChannel,results)
	}
	fmt.Println("Check unauthorized service start ...")
	//任务分发
	for _,oneTask:=range myTaskList{
		taskChannel<-oneTask
	}
	taskChannel<-common.Task{ServiceName:""}
	wg.Wait()
	fmt.Println("Unauthorized Scan finish!")
	fmt.Println(results)
	return nil
}

