package cache

import "time"

const (
	NoExpiration      time.Duration = -1          // 不会过期
	DefaultExpiration time.Duration = 0           // 默认的过期时间，在cache里面设置
	storePersisted    string        = "persisted" // 持久化存储未过期的key-value文件名前缀
	storeExpired      string        = "expired"   // 持久化存储过期的key-value文件名前缀
	SLICE             string        = "slice"     // 切片类型
	INT               string        = "int"       // int int8 int16 int32 int64
	UINT              string        = "uint"      // uint uint8 uint16 uint32 uint64
	MAP               string        = "map"       // map类型
	FLOAT             string        = "float"     // float32 float64
	CUSTOM            string        = "custom"    // 用户自定义的数据类型
)
