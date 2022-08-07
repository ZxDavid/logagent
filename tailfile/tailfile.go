package tailfi

import (
	"logagent/common"
	"logagent/kafka"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

//tail相关

type tailTask struct {
	path  string
	topic string
	tObj  *tail.Tail
}

func newTailTask(path, topic string) *tailTask {

	tt := &tailTask{
		path:  path,
		topic: topic,
	}
	return tt
}

func (t *tailTask) Init() (err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.tObj, err = tail.TailFile(t.path, config)
	return
}

func (t *tailTask) run() {
	//读取日志，发往kafka
	logrus.Infof("collect for path:%s is running...", t.path)
	for {
		//循环读数据
		line, ok := <-t.tObj.Lines

		if !ok {
			logrus.Warn("tail file close reopen, filename:%s\n", t.path)
			time.Sleep(time.Second)
			continue
		}
		//如果是空行就略过
		if len(strings.Trim(line.Text, "\r")) == 0 {
			continue
		}
		//利用通道将同步的代码改为异步
		//把读出来的一行日志包装成kafka里面的msg类型，丢到通道中
		msg := &sarama.ProducerMessage{}
		msg.Topic = t.topic
		msg.Value = sarama.StringEncoder(line.Text)
		kafka.MsgChan(msg)
	}

}

func Init(allConf []common.CollectEntry) (err error) {
	//allConf里面存了若干个日志的收集项
	//针对每一个日志收集项创建一个对应的tailobj

	for _, conf := range allConf {
		tt := newTailTask(conf.Path, conf.Topic) //创建一个日志收集的任务
		err = tt.Init()                          //打开
		if err != nil {
			logrus.Errorf("create tailObj for path:%s failed, err:%v", conf.Path, err)
			continue
		}
		logrus.Infof("create a tail task for path:%s success.", conf.Path)
		//去收集日志
		go tt.run()
	}
	return
}
