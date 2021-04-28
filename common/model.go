package common

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dutchcoders/goftp"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strings"
	"time"
)

var (
	RedisCon       *redis.Client
	FtpCon         *goftp.FTP
	MongoDBCon     *mongo.Client
	MysqlDBCon        *gorm.DB
	MssqlDBCon        *gorm.DB

	ftpNameSlice   =[]string{"anonymous","ftp"}    //空密码账户
	ftpPassSlice=[]string{"","anonymous@gmail.com"}
	mysqlNameSlice =[]string{"mysql"}
	mssqlNameSlice =[]string{"sa"}

	ScanFuncMap map[string]ScanFunc
)
type Task struct {
	ServiceName string
	IpPort string
	IsBrute bool
}
type Service struct {
	Service string `json:"service"`
	IpPortList []string `json:"ip_port_list"`
}
type ServiceList struct {
	ServiceList []Service `json:"service_list"`
}
type Result struct {   //统一输出结果
	IpPort string      //192.168.1.1:3306
	UserPass string     //admin/pass
}
type ScanFunc func(ipPort string)(userPass string,isSuccess bool)

func init()  {
	ScanFuncMap["mongoDB"]=CheckMongoDB
	ScanFuncMap["redis"]=CheckRedis
	ScanFuncMap["memcached"]=CheckMemCache
	ScanFuncMap["ftp"]=CheckFtp
	ScanFuncMap["elasticsearch"]=CheckElasticsearch   //9200
}
func dealWithEla(s string)bool{
	temp:=strings.Split(s,"\n")[1:]
	for _,i:=range temp{
		if i=="/_cat/master"{
			return true
		}
	}
	return false
}

//统一接收ip port
func CheckHadoop(ipPort string)(string,bool){    //50070
	return "", false
}
func CheckDocker(ipPort string)(string,bool){    //2375
	return "", false
}
func CheckCouchDB(ipPort string) (string,bool) {   //5984 443等
	return "", false
}
func CheckZooKeeper(ipPort string)(string,bool)  {  //2181
	return "",false
}
func CheckElasticsearch(ipPort string)(string,bool){
	resp,err:=http.Get("http://"+ipPort+"/_cat")
	if err!=nil{
		return "",false
	}
	temp,_:=goquery.NewDocumentFromReader(resp.Body)
	if resp.StatusCode==200 && dealWithEla(temp.Text()){
		return "/",true
	}
	return "",false
}
func CheckRedis(ipPort string) (string,bool){
	RedisCon =redis.NewClient(&redis.Options{
		Addr:     ipPort,
		Password: "", // no password set
		DB:       0,  // use default DB
		DialTimeout:3,
	})
	_,err:= RedisCon.Ping().Result()
	RedisCon.Close()
	if err!=nil{
		return "",false
	}else {
		return "/",true
	}
}
func CheckMemCache(ipPort string)(string,bool ) {
	mc:=memcache.New(ipPort)
	mc.Timeout=3*time.Second
	if mc==nil{
		return "",false
	}
	err:=mc.Ping()
	if err!=nil{
		return "",false
	}else {
		return "/",true
	}
}
func CheckFtp(ipPort string)(string,bool)  {
	var err error
	FtpCon,err=goftp.Connect(ipPort)
	if err!=nil {
		fmt.Println("Unable to reach ftp target:" + ipPort)
		return "",false
	}
	for _,userName :=range ftpNameSlice {
		for _, userPass := range ftpPassSlice {
			err = FtpCon.Login(userName, userPass)
			if err == nil {
				FtpCon.Close()
				return userName+"/"+userPass,true
			}
		}
	}
	FtpCon.Quit()
	return "",false
}

func CheckTelnet(ipPort string)(bool,string){
	return false,""
}

func CheckMongoDB(ipPort string)(string,bool)  {
	appUrl:="mongodb://"+ipPort
	clientOption:=options.Client().ApplyURI(appUrl)
	var err error
	MongoDBCon,err=mongo.Connect(context.TODO(),clientOption)
	if err!=nil{
		fmt.Println("Unable to connect MongoDB:"+ipPort)
		return "",false
	}
	err= MongoDBCon.Ping(context.TODO(),nil)
	if err!=nil{
		MongoDBCon.Disconnect(context.TODO())
		fmt.Println("Check mongoDB link failed:"+ipPort)
		return "",false
	}
	MongoDBCon.Disconnect(context.TODO())
	return "/",true
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