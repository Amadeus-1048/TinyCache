package tinyCache

import (
	"reflect"
	"testing"
)

// 定义一个函数类型 F，并且实现接口 A 的方法，然后在这个方法中调用自己。
// 这是 Go 语言中将其他函数（参数返回值定义与 F 一致）转换为接口 A 的常用技巧。
func Test_Getter(t *testing.T) {
	// 借助 GetterFunc 的类型转换，将一个匿名回调函数转换成了接口 f Getter。
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	// 调用该接口的方法 f.Get(key string)，实际上就是在调用匿名回调函数
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}
