package qConfig

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	Settings Setting
)

type Setting struct {
	Application Application `mapstructure:"application"`
	Log         Log         `mapstructure:"log"`
	Database    Database    `mapstructure:"database"`
	Redis       Redis       `mapstructure:"redis"`
	Jwt         Jwt         `mapstructure:"jwt"`
	Email       Email       `mapstructure:"email"`
	MingMang    MingMang    `mapstructure:"mingmang"`
	BDH         BDH         `mapstructure:"bdh"`
	VBee        VBee        `mapstructure:"vbee"`
	Kafka       Kafka       `mapstructure:"kafka"`
	AzureTTS    AzureTTS    `mapstructure:"azuretts"`
	US3         US3         `mapstructure:"us3"`
	TX          TX          `mapstructure:"tx"`
	CNTX        TX          `mapstructure:"cntx"`
	Task        Task        `mapstructure:"task"`
}

type Application struct {
	Host                string        `mapstructure:"host"`
	Port                int           `mapstructure:"port"`
	Name                string        `mapstructure:"name"`
	Mode                string        `mapstructure:"mode"`
	Domain              string        `mapstructure:"domain"`
	IsHttps             bool          `mapstructure:"isHttps"`
	WriteTimeoutSeconds time.Duration `mapstructure:"writeTimeoutSeconds"`
	ReadTimeoutSeconds  time.Duration `mapstructure:"readTimeoutSeconds"`
}

type Log struct {
	Compress      bool   `mapstructure:"compress"`
	ConsoleStdout bool   `mapstructure:"consoleStdout"`
	FileStdout    bool   `mapstructure:"fileStdout"`
	Level         string `mapstructure:"level"`
	LocalTime     bool   `mapstructure:"localtime"`
	Path          string `mapstructure:"path"`
	MaxSize       int    `mapstructure:"maxSize"`
	MaxAge        int    `mapstructure:"maxAge"`
	MaxBackups    int    `mapstructure:"maxBackups"`
	CollectorURL  string `mapstructure:"collectorURL"`
	Insecure      bool   `mapstructure:"insecure"`
}

type Database struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Name        string `mapstructure:"name"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	MaxOpenConn int    `mapstructure:"maxOpenConn"`
	MaxIdleConn int    `mapstructure:"maxIdleConn"`
}

type Redis struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	DataBase    int    `mapstructure:"dataBase"`
	Password    string `mapstructure:"password"`
	MaxIdleConn int    `mapstructure:"maxIdleConn"`
	MinIdleConn int    `mapstructure:"minIdleConn"`
}

type Jwt struct {
	PrivateKey     string        `mapstructure:"privateKey"`
	PublicKey      string        `mapstructure:"publicKey"`
	Secret         string        `mapstructure:"secret"`
	TimeoutSeconds time.Duration `mapstructure:"timeoutSeconds"`
}

type Email struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Alias    string `mapstructure:"alias"`
	Password string `mapstructure:"password"`
}

type MingMang struct {
	AppId   string `mapstructure:"appId"`
	UserKey string `mapstructure:"userKey"`
}

type BDH struct {
	AppId      string `mapstructure:"appId"`
	AppKey     string `mapstructure:"userKey"`
	WSDomain   string `mapstructure:"wsDomain"`
	HTTPDomain string `mapstructure:"httpDomain"`
}

type VBee struct {
	AppId string `mapstructure:"appId"`
	Token string `mapstructure:"token"`
}

type Kafka struct {
	Brokers string `mapstructure:"brokers"`
	Topics  string `mapstructure:"topics"`
	Group   string `mapstructure:"group"`
	Version string `mapstructure:"version"`
}

type AzureTTS struct {
	URL string `mapstructure:"url"`
	Key string `mapstructure:"key"`
}

type US3 struct {
	PublicKey       string `mapstructure:"publicKey"`
	PrivateKey      string `mapstructure:"privateKey"`
	BucketHost      string `mapstructure:"bucketHost"`
	FileHost        string `mapstructure:"fileHost"`
	BucketName      string `mapstructure:"bucketName"`
	VerifyUploadMD5 bool   `mapstructure:"VerifyUploadMD5"`
}

type TX struct {
	Appkey               string `mapstructure:"appkey"`
	AccessToken          string `mapstructure:"accessToken"`
	TxRequestTimeout     uint32 `mapstructure:"txRequestTimeout"`
	MaxConcurrentSession int    `mapstructure:"maxConcurrentSession"`
	HTTPDomain           string `mapstructure:"HTTPDomain"`
	WSSDomain            string `mapstructure:"WSSDomain"`
	PushKey              string `mapstructure:"pushKey"`
	PushURL              string `mapstructure:"pushURL"`
	PullURL              string `mapstructure:"pullURL"`
}

type Task struct {
	Uri string `mapstructure:"uri"`
}

// 载入配置文件
func Setup(path string) {
	viper.SetConfigFile(path)
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Read config file fail: %s", err.Error())
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatalf("Parse config file fail: %s", err.Error())
	}
	err = viper.Unmarshal(&Settings)
	if err != nil {
		log.Fatalf("Unmarshal config fail: %s", err.Error())
	}

}
