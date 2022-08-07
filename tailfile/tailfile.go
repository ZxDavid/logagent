package tailfi

import (
	"logagent/etcd"

	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

//tail相关

var (
	TailObj []tail.Tail
)

func Init() (err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	for _, value := range etcd.EtcdConf {
		Ta, err := tail.TailFile(value.Path, config)
		TailObj = append(TailObj, *Ta)
		if err != nil {
			logrus.Error("tailfile: create tailObj for path:%s failed,err:%v\n", value.Path, err)
			continue
		}
	}

	//打开文件开始读取数据
	// TailObj, err = tail.TailFile(filename, config)
	// if err != nil {
	// 	logrus.Error("tailfile: create tailObj for path:%s failed,err:%v\n", filename, err)
	// 	return
	// }

	return
}
