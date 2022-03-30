package cache

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Cache struct {
	defaultExpiration time.Duration      // 默认过期时间
	items             map[string]Item    // 存放数据
	delMap            map[string]delItem // 存放被删除的数据
	prefixTree        *trie              // 提供key的前缀查询
	mu                sync.RWMutex       // 读写锁
	size              int                // 记录当前的cache中key的数量
	gc                *garcoll           // 自动清理过期的key
	persistSeq        int                // 持久化文件的序号
}

// NewClient 新建一个Cache客户端，需要传入的参数是默认的到期时间和过期清理周期
func NewClient(expiredTime time.Duration, cleanupInterval time.Duration) *Cache {
	cache := &Cache{
		defaultExpiration: expiredTime,
		items:             make(map[string]Item),
		delMap:            make(map[string]delItem),
		prefixTree:        newTrie(),
		mu:                sync.RWMutex{},
		size:              0,
		persistSeq:        1,
		gc: &garcoll{
			interval: cleanupInterval,
			stop:     make(chan bool),
		},
	}

	go cache.gc.Run(cache)
	return cache
}

// Set 加入一个新的key-value或者更新旧的key-value
func (c *Cache) Set(k string, x interface{}, d time.Duration) {

	// 获取过期时间
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	// 写入
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[k] = Item{
		itemType:   c.getType(x),
		Object:     x,
		Expiration: e,
	}
	c.size++
	c.insertKey(k)
}

// SetDefault 使用默认的过期时间写入，不用传入过期时间
func (c *Cache) SetDefault(k string, x interface{}) {
	c.Set(k, x, DefaultExpiration)
}

// Add 只有当key不存在或者过期时才可以加入
func (c *Cache) Add(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.get(k); ok {
		return fmt.Errorf("Item %s already exists", k)
	}

	c.set(k, x, d)
	c.insertKey(k)
	c.size++
	return nil
}

// Replace 只有当key存在且未过期的时候可以调用，替换新的value
func (c *Cache) Replace(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.get(k); !ok {
		return fmt.Errorf("Item %s doesn't exist", k)
	}

	c.set(k, x, d)
	return nil
}

// Get 获取指定key对应的value
func (c *Cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[k]
	if !ok {
		return nil, false
	}

	if c.items[k].expired() {
		c.autoDelete(k)
		return nil, false
	}
	return item.Object, true
}

// GetWithExpiration 获取指定key对应的value和过期时间
func (c *Cache) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 不存在这个key
	item, ok := c.items[k]
	if !ok {
		return nil, time.Time{}, false
	}

	// 过期
	if c.items[k].expired() {
		c.autoDelete(k)
		return nil, time.Time{}, false
	}

	// todo: 把int64转为时间 应该是还剩多少时间过期
	return item.Object, time.Time{}, true
}

// Delete 删除指定的key
func (c *Cache) Delete(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 不存在这个key，其实也应该返回true
	if _, ok := c.items[k]; !ok {
		return
	}

	c.manualDelete(k)
}

// Increment 为指定的key增加n，
// key对应的value必须是一个数字类型
func (c *Cache) Increment(k string, n interface{}) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 当前key不存在
	val, ok := c.items[k]
	if !ok || val.expired() {
		return nil, fmt.Errorf("item %s not found", k)
	}

	switch n.(type) {
	case int:
		return c.incrementInt(k, n)
	case int8:
		return c.incrementInt8(k, n)
	case int16:
		return c.incrementInt16(k, n)
	case int32:
		return c.incrementInt32(k, n)
	case int64:
		return c.incrementInt64(k, n)
	case uint:
		return c.incrementUint(k, n)
	case uint8:
		return c.incrementUint8(k, n)
	case uint16:
		return c.incrementUint16(k, n)
	case uint32:
		return c.incrementUint32(k, n)
	case uint64:
		return c.incrementUint64(k, n)
	case float32:
		return c.incrementFloat32(k, n)
	case float64:
		return c.incrementFloat64(k, n)
	default:
		return nil, fmt.Errorf("the value for %s can not be increased", k)
	}
}

//
func (c *Cache) Decrement(k string, n interface{}) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.items[k]
	if !ok || val.expired() {
		c.autoDelete(k)
		return nil, fmt.Errorf("item %s not found", k)
	}

	switch n.(type) {
	case int:
		return c.decrementInt(k, n)
	case int8:
		return c.decrementInt8(k, n)
	case int16:
		return c.decrementInt16(k, n)
	case int32:
		return c.decrementInt32(k, n)
	case int64:
		return c.decrementInt64(k, n)
	case uint:
		return c.decrementUint(k, n)
	case uint8:
		return c.decrementUint8(k, n)
	case uint16:
		return c.decrementUint16(k, n)
	case uint32:
		return c.decrementUint32(k, n)
	case uint64:
		return c.decrementUint64(k, n)
	case float32:
		return c.decrementFloat32(k, n)
	case float64:
		return c.decrementFloat64(k, n)
	default:
		return nil, fmt.Errorf("the value for %s can not be decreased", k)
	}
}

// 关于map的操作
// todo:

// IsExistedKey 查询某个key是否存在
func (c *Cache) IsExistedKey(k string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if _, ok := c.items[k]; ok {
		return true
	}
	return false
}

// IsExistedKeyWithPrefix 查询某个前缀是否存在
func (c *Cache) IsExistedKeyWithPrefix(prefix string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.searchWithPrefix(prefix)
}

// Persist 持久化缓存数据到磁盘，包括有效数据和被删除的数据
func (c *Cache) Persist() error {
	c.mu.Lock()
	itemDir := storePersisted + strconv.Itoa(c.persistSeq)
	delItemDir := storeExpired + strconv.Itoa(c.persistSeq)
	c.persistSeq++
	c.mu.Unlock()

	// 持久化数据
	fp1, err := os.Create(itemDir)
	if err != nil {
		return err
	}
	err = c.saveItem(fp1)
	if err != nil {
		return err
	}
	fp1.Close()

	// 持久化删除的数据
	fp2, err := os.Create(delItemDir)
	if err != nil {
		return err
	}
	err = c.saveDel(fp2)
	if err != nil {
		return err
	}
	fp2.Close()
	return nil
}

// Load 加载最新的有效数据文件
// todo: 怎么找到最新的文件？
// 读取当前目录所有文件的名称，序号最大的就是最新的，同时将自己的序号更新+1
func (c *Cache) Load(seq int) error {
	itemDir := storePersisted + strconv.Itoa(seq)
	fp, err := os.Open(itemDir)
	if err != nil {
		return err
	}
	defer fp.Close()

	return c.load(fp, seq)
}

// Items 复制所有未过期的items，生成一个新的map并返回
func (c *Cache) Items() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()

	m := make(map[string]Item, len(c.items))
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			continue
		}
		m[k] = v
	}
	return m
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.size
}

// Flush 清空当前数据
func (c *Cache) Flush() {
	c.Persist()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = map[string]Item{}
}

func (c *Cache) StopGC() {
	c.gc.stop <- true
}

// for test
func (c *Cache) SearchDel(k string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.delMap[k]; ok {
		return true
	}
	return false
}
