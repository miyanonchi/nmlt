package main

import (
  "testing"

  "github.com/yuin/gopher-lua"
)

var L = lua.NewState()

func SetValueWrap(t *testing.T, key string, value string) {
  // 引数をセット
  L.Push(lua.LString(key))
  L.Push(lua.LString(value))

  // テスト対象
  SetValue(L)
}

func AssertEqual(t *testing.T, keys []string, value string) {
  var curr interface{} = listHash
  for i, k := range(keys) {
    if _, ok := curr.(map[string]interface{})[k]; ok {
      curr = curr.(map[string]interface{})[k]

      if i + 1 == len(keys) {
        // 値が一致するかどうか
        if curr.(string) != value {
          t.Fatalf("assert fail: " + curr.(string) + " != " + value)
        }
      } else {
        curr = curr.(map[string]interface{})
      }
    } else {
      // キーがない
      t.Fatalf("assert fail: no key " + k)
    }
  }
}

func AssertPanic(t *testing.T, f func(interface{}), args ...interface{}) {
  defer func() {
    if r:= recover(); r == nil {
      t.Fatalf("did not panic.")
    }
  }()

  f(args)
}

func TestSetValue(t *testing.T) {
  /**************/
  /* 正常ケース */
  /**************/

  SetValueWrap(t, "$.key1.key2", "value1")
  AssertEqual(t, []string{"key1", "key2"}, "value1")

  /**************/
  /* 異常ケース */
  /**************/
  AssertPanic(t, SetValueWrap, t, "key1.array1", "value2")
}
