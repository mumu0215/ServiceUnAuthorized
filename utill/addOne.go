package utill

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
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
	var myService ServiceList
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
	return nil
}

