package common

import (
	"context"
	"fmt"
	"github.com/dutchcoders/goftp"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"sync"
)

var (
	RedisCon       *redis.Client
	FtpCon         *goftp.FTP
	MongoDBCon     *mongo.Client
	MysqlDBCon        *gorm.DB
	MssqlDBCon        *gorm.DB

	ftpNameSlice   =[]string{"anonymous","ftp"}    //空密码账户
	mysqlNameSlice =[]string{"mysql"}
	mssqlNameSlice =[]string{"sa"}
)

func MainWorker(wg *sync.WaitGroup,task chan Task)  {
	for tempTask :=range task{
	}
}

func worker(task chan Task)  string{
	temp:=<-task
	switch temp.ServiceName {
	case "mongoDB":
	case "redis":
	case "mysql":
	case "mssql":
	case "ssh":
	case "telnet":
	case "ftp":
	case "oracle-tns":
	default:

	}
	return ""
}


//统一接收ip port
func CheckRedis(ipPort string) bool{
	RedisCon =redis.NewClient(&redis.Options{
		Addr:     ipPort,
		Password: "", // no password set
		DB:       0,  // use default DB
		DialTimeout:3,
	})
	_,err:= RedisCon.Ping().Result()
	RedisCon.Close()
	if err!=nil{
		return false
	}else {
		return true
	}
}
func CheckFtp(ipPort string)(bool,string)  {
	var err error
	FtpCon,err=goftp.Connect(ipPort)
	if err!=nil {
		fmt.Println("Unable to reach ftp target:" + ipPort)
		return false,""
	}
	for _,userName :=range ftpNameSlice{
		err= FtpCon.Login(userName,"")
		if err==nil{
			FtpCon.Close()
			return true,userName
		}
	}
	FtpCon.Quit()
	return false,""
}

func CheckTelnet(ipPort string)(bool,string){
	return false,""
}

func CheckMongoDB(ipPort string)bool  {
	appUrl:="mongodb://"+ipPort
	clientOption:=options.Client().ApplyURI(appUrl)
	var err error
	MongoDBCon,err=mongo.Connect(context.TODO(),clientOption)
	if err!=nil{
		fmt.Println("Unable to connect MongoDB:"+ipPort)
		return false
	}
	err= MongoDBCon.Ping(context.TODO(),nil)
	if err!=nil{
		MongoDBCon.Disconnect(context.TODO())
		fmt.Println("Check mongoDB link failed:"+ipPort)
		return false
	}
	MongoDBCon.Disconnect(context.TODO())
	return true
}

func CheckMysql(ipPort string)(bool,string)  {
	var err error
	connectStr:=mysqlNameSlice[0]+":@tcp("+ipPort+")/test?charset=utf8&parseTime=True&loc=Local&timeout=4s"
	MysqlDBCon,err=gorm.Open("mysql",connectStr)
	if err!=nil{
		fmt.Println("Unable to connect Mysql:"+ipPort)
		return false, ""
	}
	err=MysqlDBCon.DB().Ping()
	if err!=nil{
		fmt.Println("Error in ping Mysql connect:"+ipPort)
		return false, ""
	}
	return true,mysqlNameSlice[0]
}
func CheckMssql(ipPort string)(bool,string){
	temp:=strings.Split(ipPort,":")
	var err error
	connectStr:=fmt.Sprintf("server=%s;port=%s;database=master;user id=%s;password=%s;timeout=4s",
		temp[0],temp[1],mssqlNameSlice[0],"")
	MssqlDBCon,err=gorm.Open("mssql",connectStr)
	if err!=nil{
		fmt.Println("Unable to connect Mssql:"+ipPort)
		return false, ""
	}
	if err=MssqlDBCon.DB().Ping();err!=nil{
		fmt.Println("Error in ping Mssql connect:"+ipPort)
		return false, ""
	}
	return true,mssqlNameSlice[0]
}
