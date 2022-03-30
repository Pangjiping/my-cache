package cache

import (
	"reflect"
	"time"
)

// Item 构建一个存储对象结构
type Item struct {
	itemType   string      // 存储对象的类型
	Object     interface{} // 存储对象
	Expiration int64       // 过期时间
}

// expired 判断当前的数据是否过期，需要在外部加读锁
func (item Item) expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

// getType 获取当前item的数据类型
func (item Item) getType() string {
	return reflect.TypeOf(item.Object).String()
}

// delItem 被删除的数据会留一个备份
type delItem struct {
	itemType      string      // 数据类型
	Object        interface{} // 数据内容
	isAutoCleanup bool        // 是否是自动清理的
	isExpired     bool        // 在删除它的时候它是不是过期的
	deletedAt     time.Time   // 删除的时间点
}
