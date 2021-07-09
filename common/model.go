package common

import (
	"context"
	"crypto/tls"
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
	"net"
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
	Client *http.Client

	ftpNameSlice   =[]string{"anonymous","ftp"}    //空密码账户
	ftpPassSlice=[]string{"","anonymous@gmail.com"}
	mysqlNameSlice =[]string{"mysql"}
	mssqlNameSlice =[]string{"sa"}
	author=[]string{"Basic Z3Vlc3Q6Z3Vlc3Q=","Basic cm9vdDpyb290"}
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
	UserPass string     //admin/pass提示字符串
}
type ScanFunc func(ipPort string)(userPass string,isSuccess bool)  //userPass可以填充提示信息字符串

func init()  {
	Client=&http.Client{
		Timeout:time.Duration(5)*time.Second,
		Transport: &http.Transport{
			//参数未知影响，目前不使用
			//TLSHandshakeTimeout: time.Duration(timeout) * time.Second,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},}
	ScanFuncMap["mongoDB"]=CheckMongoDB
	ScanFuncMap["redis"]=CheckRedis
	ScanFuncMap["memcached"]=CheckMemCache
	ScanFuncMap["ftp"]=CheckFtp
	ScanFuncMap["zookeeper"]=CheckZooKeeper

	ScanFuncMap["elasticsearch"]=CheckElasticsearch   //9200
	ScanFuncMap["hadoop"]=CheckHadoop
	ScanFuncMap["docker"]=CheckDocker
	ScanFuncMap["couchdb"]=CheckCouchDB
	ScanFuncMap["jenkins"]=CheckJenkins
	ScanFuncMap["kibana"]=CheckKibana
	ScanFuncMap["kubernetes"]=CheckKubernetes
	ScanFuncMap["activemq"]=CheckActiveMQ
	ScanFuncMap["rabbitmq"]=CheckRabbitMQ
	ScanFuncMap["sprintboot"]=CheckSpringBootActuator
	ScanFuncMap["jboss"]=CheckJBoss
	ScanFuncMap["dubbo"]=CheckDubbo
	ScanFuncMap["atlassiancrowd"]=CheckAtlassianCrowd
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

//统一接收ip port
func CheckHadoop(ipPort string)(string,bool){    //50070 /ws/v1/cluster/info
	statusCode,temp,isSend:=sendGetRequest("http://"+ipPort+"/ws/v1/cluster/info")
	if !isSend{
		return "", false
	}
	if statusCode==200&&strings.Contains(temp,"resourceManagerVersionBuiltOn")&&strings.Contains(temp,"hadoopVersion"){
		return ipPort+" has Hadoop unAuthorized",true
	}else {
		return "", false
	}
}
func CheckDocker(ipPort string)(string,bool){    //2375
	statusCode,tempStr,isSend:=sendGetRequest("http://"+ipPort+"/info")
	if !isSend{
		return "", false
	}
	if statusCode==200&&strings.Contains(tempStr,"KernelVersion")&&strings.Contains(tempStr,"RegistryConfig")&&strings.Contains(tempStr,"DockerRootDir"){
		return ipPort+" has Docker api unAuthorized", true
	}else {
		return "", false
	}
}
func CheckCouchDB(ipPort string) (string,bool) {   //5984 6984 443等
	statusCode,tempStr,isSend:=sendGetRequest("http://"+ipPort+"/_config")
	if !isSend{
		return "", false
	}
	if statusCode==200&&strings.Contains(tempStr,"httpd_design_handlers")&&strings.Contains(tempStr,"external_manager")&&strings.Contains(tempStr,"replicator_manager"){
		return ipPort+" has CouchDB unAuthorized", true
	}else {
		return "", false
	}
}
func CheckZooKeeper(ipPort string)(string,bool)  {  //2181   service端口
	conn,err:=net.DialTimeout("tcp",ipPort,time.Second*5)
	if err!=nil{
		return "", false
	}
	conn.Write([]byte("envi"))
	d:=make([]byte,4096)
	in,err:=conn.Read(d)
	if err!=nil{
		return "", false
	}
	if strings.HasPrefix(string(d[:in]),"Environment:"){
		return ipPort+" has Zookeeper unAuthorized", true
	}else {
		return "", false
	}
}
func CheckJenkins(ipPort string)(string,bool){   //web port
	statusCode,_,isSend:=sendGetRequest("http://"+ipPort+"/script")
	if !isSend{
		return "", false
	}
	if statusCode==200{
		return ipPort+" has Jenkins unAuthorized", true
	}
	return "", false
}
func CheckVNC(ipPort string)(string,bool){    //5900 5901
	return "", false
}
func CheckNfs(ipPort string)(string,bool)  {    //2049
	return "", false
}
func CheckRsync(ipPort string)(string,bool){    //873
	return "", false
}
func CheckKibana(ipPort string)(string,bool){  //5601 web
	statusCode,tempStr,isSend:=sendGetRequest("http://"+ipPort+"/app/kibana")
	if !isSend{
		return "", false
	}
	if statusCode==200&&strings.Contains(tempStr,".kibanaWelcomeView"){
		return ipPort+" has Kibana unAuthorized", true
	}
	return "", false
}
func CheckKubernetes(ipPort string)(string,bool){   //web
	statusCode,tempStr,isSend:=sendGetRequest("http://"+ipPort)
	if !isSend{
		return "", false
	}
	statusCode1,_,isSend:=sendGetRequest("http://"+ipPort+"/api")
	if !isSend{
		return "",false
	}
	if statusCode==200&&strings.Contains(tempStr,"/api")&&statusCode1==200{
		return ipPort+" has k8s unAuthorized", true
	}
	return "", false
}
func CheckActiveMQ(ipPort string)(string,bool){    //8161
	statusCode,tempStr,isSend:=sendGetRequest("http://"+ipPort+"/admin")
	if !isSend{
		return "", false
	}
	if statusCode==200&&strings.Contains(tempStr,"ActiveMQ Console"){
		return ipPort+" has ActiveMQ unAuthorized", true
	}else if statusCode==401{
		req,err:=http.NewRequest("GET","http://"+ipPort+"/admin",nil)
		if err!=nil{
			return "", false
		}
		req.Header.Add("Authorization","Basic YWRtaW46YWRtaW4=")
		req.Header.Add("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:78.0) Gecko/20100101 Firefox/78.0")
		respond,err:=Client.Do(req)
		if err!=nil{
			return "", false
		}
		defer respond.Body.Close()
		temp,err:=goquery.NewDocumentFromReader(respond.Body)
		if respond.StatusCode==200 && strings.Contains(temp.Text(),"ActiveMQ Console"){
			return ipPort+" has ActiveMQ unAuthorized, admin/admin", true
		}
	}
	return "", false
}
func CheckRabbitMQ(ipPort string)(string,bool){   //5672,15672（guest/guest）
	req,err:=http.NewRequest("GET","http://"+ipPort+"/api/whoami",nil)
	if err!=nil{
		return "", false
	}
	req.Header.Add("authorization","Basic Z3Vlc3Q6Z3Vlc3Q=")
	req.Header.Add("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:78.0) Gecko/20100101 Firefox/78.0")
	respond,err:=Client.Do(req)
	if err!=nil{
		return "", false
	}
	defer respond.Body.Close()
	temp,err:=goquery.NewDocumentFromReader(respond.Body)
	if respond.StatusCode==200 && strings.Contains(temp.Text(),"\"name\":\"guest"){
		return ipPort+" has RabbitMQ unAuthorized, guest/guest", true
	}
	return "", false
}
func CheckSpringBootActuator(ipPort string)(string,bool){    //web
	statusCode1,temp1,isSend:=sendGetRequest("http://"+ipPort+"/env")
	statusCode2,temp2,isSend1:=sendGetRequest("http://"+ipPort+"/actuator/env")
	if isSend{
		if statusCode1==200&&strings.Contains(temp1,"java.version")&&strings.Contains(temp1,"os.arch"){
			return ipPort+" has SpringBoot Actuator unAuthorized in /env", true
		}
	}else if isSend1{
		if statusCode2==200&&strings.Contains(temp2,"java.version")&&strings.Contains(temp2,"os.arch"){
			return ipPort+" has SpringBoot Actuator unAuthorized in /actuator/env", true
		}
	}
	return "", false
}
func CheckJBoss(ipPort string)(string,bool){     //web   /jmx-console
	statusCode,temp,isSend:=sendGetRequest("http://"+ipPort+"/jmx-console/")
	if isSend{
		if statusCode==200&&strings.Contains(temp,"jboss.management.local")&&strings.Contains(temp,"jboss.web"){
			return ipPort+" has JBoss unAuthorized", true
		}
	}
	return "", false
}
func CheckDubbo(ipPort string)(string,bool){     //存在web和协议端口两种测试方式 Dubbo Admin
	//Authorization: Basic Z3Vlc3Q6Z3Vlc3Q=
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
			}else if index==1{
				return ipPort+" has Dubbo unAuthorized, root/root", true
			}
		}
	}
	return "", false
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

func CheckElasticsearch(ipPort string)(string,bool){
	statusCode,temp,isSend:=sendGetRequest("http://"+ipPort+"/_cat")
	if !isSend{
		return "", false
	}
	if statusCode==200 && dealWithEla(temp){
		return ipPort+" has Elasticsearch unAuthorized",true
	}
	return "",false
}
func CheckRedis(ipPort string) (string,bool){
	RedisCon =redis.NewClient(&redis.Options{
		Addr:     ipPort,
		Password: "", // no password set
		DB:       0,  // use default DB
		DialTimeout:time.Duration(4),
	})
	_,err:= RedisCon.Ping().Result()
	RedisCon.Close()
	if err!=nil{
		return "",false
	}else {
		return ipPort+" has Redis unAuthorized",true
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
		return ipPort+" has MemCache unAuthorized",true
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
				return ipPort+" has Ftp unAuthorized: "+userName+"/"+userPass,true
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
		return "",false
	}
	err= MongoDBCon.Ping(context.TODO(),nil)
	if err!=nil{
		MongoDBCon.Disconnect(context.TODO())
		return "",false
	}
	MongoDBCon.Disconnect(context.TODO())
	return ipPort+" has MongoDB unAuthorized",true
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