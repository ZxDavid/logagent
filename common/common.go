package common

//CollectEntry要收集的日志的结构体的配置项
type CollectEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}
