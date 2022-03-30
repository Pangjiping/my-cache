package cache

import (
	"encoding/gob"
	"fmt"
	"io"
	"reflect"
	"time"
)

// getType 获取数据类型
func (c *Cache) getType(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

// autoDelete 自动删除某个过期key，这里一定是过期的
func (c *Cache) autoDelete(k string) {
	item := c.items[k]
	del := delItem{
		itemType:      item.getType(),
		Object:        item.Object,
		isAutoCleanup: true,
		isExpired:     true,
		deletedAt:     time.Now(),
	}

	delete(c.items, k)
	c.delMap[k] = del
}

// manualDelete 手动删除一个key
func (c *Cache) manualDelete(k string) {
	item := c.items[k]
	del := delItem{
		itemType:      item.getType(),
		Object:        item.Object,
		isAutoCleanup: false,
		isExpired:     false,
		deletedAt:     time.Now(),
	}

	// 先判断这个key有没有过期
	if c.items[k].expired() {
		del.isExpired = true
	}

	delete(c.items, k)
	c.delMap[k] = del
}

// delete 扫描所有key，过期删除
func (c *Cache) delete() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, item := range c.items {
		if !item.expired() {
			continue
		}
		delete(c.items, k)
		c.delMap[k] = delItem{
			itemType:      item.getType(),
			Object:        item.Object,
			isAutoCleanup: true,
			isExpired:     true,
			deletedAt:     time.Now(),
		}
	}
}

func (c *Cache) incrementInt(k string, n interface{}) (int, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case int:
		item.Object = item.Object.(int) + n.(int)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(int), nil
}

func (c *Cache) incrementInt8(k string, n interface{}) (int8, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case int8:
		item.Object = item.Object.(int8) + n.(int8)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(int8), nil
}

func (c *Cache) incrementInt16(k string, n interface{}) (int16, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case int16:
		item.Object = item.Object.(int16) + n.(int16)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(int16), nil
}

func (c *Cache) incrementInt32(k string, n interface{}) (int32, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case int32:
		item.Object = item.Object.(int32) + n.(int32)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(int32), nil
}

func (c *Cache) incrementInt64(k string, n interface{}) (int64, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case int64:
		item.Object = item.Object.(int64) + n.(int64)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(int64), nil
}

func (c *Cache) incrementUint(k string, n interface{}) (uint, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case uint:
		item.Object = item.Object.(uint) + n.(uint)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(uint), nil
}

func (c *Cache) incrementUint8(k string, n interface{}) (uint8, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case uint8:
		item.Object = item.Object.(uint8) + n.(uint8)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(uint8), nil
}

func (c *Cache) incrementUint16(k string, n interface{}) (uint16, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case uint16:
		item.Object = item.Object.(uint16) + n.(uint16)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(uint16), nil
}

func (c *Cache) incrementUint32(k string, n interface{}) (uint32, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case uint32:
		item.Object = item.Object.(uint32) + n.(uint32)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(uint32), nil
}

func (c *Cache) incrementUint64(k string, n interface{}) (uint64, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case uint64:
		item.Object = item.Object.(uint64) + n.(uint64)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(uint64), nil
}

func (c *Cache) incrementFloat32(k string, n interface{}) (float32, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case float32:
		item.Object = item.Object.(float32) + n.(float32)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(float32), nil
}

func (c *Cache) incrementFloat64(k string, n interface{}) (float64, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case float64:
		item.Object = item.Object.(float64) + n.(float64)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(float64), nil
}

func (c *Cache) decrementInt(k string, n interface{}) (int, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case int:
		item.Object = item.Object.(int) - n.(int)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(int), nil
}

func (c *Cache) decrementInt8(k string, n interface{}) (int8, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case int8:
		item.Object = item.Object.(int8) - n.(int8)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(int8), nil
}

func (c *Cache) decrementInt16(k string, n interface{}) (int16, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case int16:
		item.Object = item.Object.(int16) - n.(int16)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(int16), nil
}

func (c *Cache) decrementInt32(k string, n interface{}) (int32, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case int32:
		item.Object = item.Object.(int32) - n.(int32)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(int32), nil
}

func (c *Cache) decrementInt64(k string, n interface{}) (int64, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case int64:
		item.Object = item.Object.(int64) - n.(int64)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(int64), nil
}

func (c *Cache) decrementUint(k string, n interface{}) (uint, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case uint:
		item.Object = item.Object.(uint) - n.(uint)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(uint), nil
}

func (c *Cache) decrementUint8(k string, n interface{}) (uint8, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case uint8:
		item.Object = item.Object.(uint8) - n.(uint8)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(uint8), nil
}

func (c *Cache) decrementUint16(k string, n interface{}) (uint16, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case uint16:
		item.Object = item.Object.(uint16) - n.(uint16)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(uint16), nil
}

func (c *Cache) decrementUint32(k string, n interface{}) (uint32, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case uint32:
		item.Object = item.Object.(uint32) - n.(uint32)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(uint32), nil
}

func (c *Cache) decrementUint64(k string, n interface{}) (uint64, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case uint64:
		item.Object = item.Object.(uint64) - n.(uint64)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(uint64), nil
}

func (c *Cache) decrementFloat32(k string, n interface{}) (float32, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case float32:
		item.Object = item.Object.(float32) - n.(float32)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(float32), nil
}

func (c *Cache) decrementFloat64(k string, n interface{}) (float64, error) {
	item := c.items[k]
	switch item.Object.(type) {
	case float64:
		item.Object = item.Object.(float64) - n.(float64)
	default:
		return 0, fmt.Errorf("type mismatch")
	}
	c.items[k] = item
	return item.Object.(float64), nil
}

func (c *Cache) set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.items[k] = Item{
		itemType:   reflect.TypeOf(x).String(),
		Object:     x,
		Expiration: e,
	}
}

func (c *Cache) get(k string) (interface{}, bool) {
	item, ok := c.items[k]
	if !ok {
		return nil, false
	}
	if item.expired() {
		c.autoDelete(k)
		return nil, false
	}
	return item.Object, true
}

// Save 使用gob编码将cache内容写到io.Writer
func (c *Cache) saveItem(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error registering item types with Gob library")
		}
	}()

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, v := range c.items {
		gob.Register(v.Object)
	}
	err = enc.Encode(&c.items)
	return
}

func (c *Cache) saveDel(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error registering item type with Gob library")
		}
	}()

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, v := range c.delMap {
		gob.Register(v.Object)
	}
	err = enc.Encode(&c.delMap)
	return
}

func (c *Cache) load(r io.Reader, seq int) error {
	dec := gob.NewDecoder(r)
	items := map[string]Item{}
	err := dec.Decode(&items)
	if err == nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		for k, v := range items {
			if ov, found := c.items[k]; !found || ov.expired() {
				c.items[k] = v
			}
		}
		c.persistSeq = seq + 1
	}
	return err
}

// insertKey 向字典树中加入key，外部加锁
func (c *Cache) insertKey(k string) {
	c.prefixTree.insert(k)
}

// 查询是否有带有某个前缀的key
func (c *Cache) searchWithPrefix(prefix string) bool {
	return c.prefixTree.startsWithPrefix(prefix)
}
