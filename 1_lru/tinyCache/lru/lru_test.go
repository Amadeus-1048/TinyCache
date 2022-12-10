package lru

import (
	"reflect"
	"testing"
)

type myString string

func (m myString) Len() int {
	return len(m)
}

func Test_Get(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", myString("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(myString)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func Test_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	Cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(Cap), nil)
	lru.Add(k1, myString(v1))
	lru.Add(k2, myString(v2))
	lru.Add(k3, myString(v3))
	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("RemoveOldest key1 failed")
	}
}

func Test_OnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("k1", myString("v1"))
	lru.Add("k2", myString("v2"))
	lru.Add("k3", myString("v3"))
	lru.Add("k4", myString("v4"))

	expect := []string{"k1", "k2", "k3", "k4"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
