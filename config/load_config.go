package config

import (
	"fmt"
	"github.com/jinzhu/configor"
	"baselib/lib-flag"
	"baselib/logger"
	"os"
	"strings"
)

type AppConfig struct {
	//Redis相关配置
	RedisDb struct {
		Addr     string
		Password string
		User     string
		Database int
	}

	//GRpc server端配置
	RpcServer struct {
		Protocol      string
		Port          string
		InterfaceName string
	}

	//上传文件配置
	OSS struct {
		AccessKeyId     string
		AccessKeySecret string
		BucketName      string
		Action          string
		Replace         string
	}

	//获取信息配置
	Fetch struct {
		FetchUrl       string
		FetchWeChatUrl string
	}

	//rabbitMQ配置信息
	RabbitMQ struct {
		Addr string
	}

	//调试配置
	DebugConfig struct {
		Debug bool
	}

	//Etcd配置
	Etcd struct {
		Addr string
	}

	//日志级别
	Logger struct {
		Level string
	}

	//App
	App struct {
		AppKey string
	}

	//MongoDB配置信息
	MongoDb struct {
		Hostsports string
		Dbname     string
		Userpass   string
		ReplicaSet string
		Role       string
	}

	//microservice001服务配置
	MicroService001 struct {
		Url string
	}

	//microservice002服务配置
	MicroService002 struct{
		Url string
	}
}

var conf AppConfig

func LoadConfig(file string) {
	conf = AppConfig{}
	err := configor.Load(&conf, file)
	if err != nil {
		logger.Error("Failed to find configuration ", file)
		os.Exit(1)
	} else {
		conf.OSS.Action = fmt.Sprintf(conf.OSS.Action, conf.OSS.BucketName)
		conf.OSS.Replace = fmt.Sprintf(conf.OSS.Replace, conf.OSS.BucketName)
		if !strings.EqualFold(conf.Logger.Level, "") {
			logger.SetLevel(conf.Logger.Level)
		}
	}
}

func init() {
	confPath := lib_flag.ConfPath
	LoadConfig(confPath)
	for k, v := range os.Args {
		logger.Info(fmt.Sprintf("k:%d, v:%s", k, v))
	}

	logger.Info("\r\n")
	logger.Info("loading config...... \r\n")
	logger.Info("config path:", confPath)
	logger.Info("app appKey:", conf.App.AppKey)
	logger.Info("redis addr:", conf.RedisDb.Addr)
	logger.Info(fmt.Sprintf("rpc port:%s. interface name:%s.",
		conf.RpcServer.Port, conf.RpcServer.InterfaceName))
	logger.Info("rabbitMQ addr:", conf.RabbitMQ.Addr)
	logger.Info("etcd addr:", conf.Etcd.Addr)
	logger.Info(fmt.Sprintf("debug config:%t", conf.DebugConfig.Debug))
	logger.Info(fmt.Sprintf("microservice001 config:%v", conf.MicroService001.Url))
	logger.Info(fmt.Sprintf("microservice002 config:%v", conf.MicroService002.Url))
}

func GetConf() AppConfig {
	return conf
}
