package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testClient = NewClient(time.Second*3, time.Second*5)

func TestSet(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testClient.SetDefault("testKey", 1)
		}()
	}
	testClient.Set("a", 10, time.Second*8)
	testClient.Set("b", 11, time.Second*9)
	testClient.Set("c", 12, -1)
	wg.Wait()
}

func TestGet(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, _ := testClient.Get("testKey")
			require.Equal(t, res, 1)
		}()
	}
	wg.Wait()
}

func TestAdd(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			err := testClient.Add("testKey", v, time.Second)
			require.Error(t, err)
		}(i)
	}
	wg.Wait()
}

func TestDelete(t *testing.T) {
	testClient.Delete("testKey")
	ok := testClient.IsExistedKey("testKey")
	require.Equal(t, false, ok)

	ok = testClient.SearchDel("testKey")
	require.Equal(t, true, ok)
}

func TestAutoDelete(t *testing.T) {
	time.Sleep(10 * time.Second)

	// 查询 testKey 应该查不到
	ok := testClient.IsExistedKey("testKey")
	require.Equal(t, false, ok)
	ok = testClient.IsExistedKey("a")
	require.Equal(t, false, ok)
	ok = testClient.IsExistedKey("b")
	require.Equal(t, false, ok)

	// 查询delMap 可以查到
	ok = testClient.SearchDel("testKey")
	require.Equal(t, true, ok)
	ok = testClient.SearchDel("a")
	require.Equal(t, true, ok)
	ok = testClient.SearchDel("b")
	require.Equal(t, true, ok)
}

func TestPersist(t *testing.T) {
	err := testClient.Persist()
	require.NoError(t, err)
}

func TestLoad(t *testing.T) {
	err := testClient.Load(1)
	require.NoError(t, err)

	ok := testClient.IsExistedKey("c")
	require.Equal(t, true, ok)
}
