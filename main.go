package main

import (
	"fmt"
	"logagent/kafka"
	tailfi "logagent/tailfile"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

//日志收集的客户端
//类似的开源项目还有filebeat
//目标：收集指定目录下的日志文件，发送到kafka中

//现在的技能包：
//1.往kafka发数据
//2.使用tailf读日志文件

//整个logaent的配置结构体
type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	Port     string `ini:"port"`
	Topic    string `ini:"topci"`
	ChanSize int64  `ini:"chan_size"`
}

type CollectConfig struct {
	LogfilePath string `ini:"logfile_path"`
}

//真正的业务逻辑
func run() (err error) {
	//TailObj -->log -->Client --> kafka
	for {
		//循环读数据
		line, ok := <-tailfi.TailObj.Lines
		if !ok {
			logrus.Error("tail file close reopen, filename:%s\n", tailfi.TailObj.Filename)
			time.Sleep(time.Second)
			continue
		}
		//利用通道将同步的代码改为异步
		//把读出来的一行日志包装成kafka里面的msg类型，丢到通道中
		msg := &sarama.ProducerMessage{}
		msg.Topic = "web_log"
		msg.Value = sarama.StringEncoder(line.Text)
		kafka.MsgChan <- msg
	}

}

func main() {
	var configObj = new(Config)
	// 0.读配置文件‘go-ini’
	// cfg, err := ini.Load("./conf/config.ini")
	// if err != nil {
	// 	logrus.Error("load config failed, err:%v", err)
	// 	return
	// }
	// kafkaAddr := cfg.Section("kafka").Key("address").String()
	// kafkaPort := cfg.Section("kafka").Key("port").String()
	// fmt.Println("addr:", kafkaAddr, ":", kafkaPort)

	err := ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Error("load config failed, err:%v", err)
		return
	}
	fmt.Printf("%#v\n", configObj)
	// 1.初始化（做好准备工作）
	err = kafka.Init([]string{configObj.KafkaConfig.Address + ":" + configObj.KafkaConfig.Port}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Error("init kafka failed, err:%v", err)
		return
	}
	logrus.Info("init kafka success!")

	// 2.根据配置中的日志路径使用tail去收集日志
	err = tailfi.Init(configObj.CollectConfig.LogfilePath)
	if err != nil {
		logrus.Error("init tailfile failed,err:%v", err)
		return
	}
	logrus.Info("init tailfile success!")

	// 3.把日志通过sarama发送kafka
	//TailObj -->log -->Client --> kafka
	err = run()
	if err != nil {
		logrus.Error("run failed, err:%v", err)
		return
	}

}
