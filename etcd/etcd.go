package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

//etcd相关操作

type CollectEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

var (
	client *clientv3.Client
)

var (
	EtcdConf []CollectEntry
)

func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v", err)
		return
	}
	return
}

//拉取日志收集项的函数

func GetConf(key string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd by key:%s failed, err:%v", key, err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warning("get len:0 conf from etcd by key:%s", key)
		return
	}

	ret := resp.Kvs[0]
	fmt.Println("ret.Value:", ret.Value)
	//ret.Value//json格式字符串
	err = json.Unmarshal(ret.Value, &EtcdConf)
	if err != nil {
		logrus.Error("json unmarshal failed, err:%v", err)
		return
	}
	return
}
