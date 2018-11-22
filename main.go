package main

import (
  "os"
  "fmt"
  "flag"
  "path/filepath"
  "encoding/json"

  "github.com/yuin/gopher-lua"

  "./jsonpath"
)

var listHash map[string]interface{} = make(map[string]interface{})

func SetValue(L *lua.LState) int {
  // luaで積んだキーを取得
  key    := L.ToString(1)
  target := L.ToString(2)

  // キーをjsonpathコンパイル
  compiled, err := jsonpath.Compile(key)

  if err != nil {
    panic(err)
  }

  var ptr map[string]interface{} = listHash
  var val map[string]interface{}
  var lastKey string

  if 0 < len(compiled.Steps()) && compiled.Steps()[0].Op() != "root" {
    panic("jsonpath must start with '$'")
  }

  for i, v := range(compiled.Steps()) {
    if v.Op() != "key" {
        continue
    }

    lastKey = v.Key()

    if i + 1 == len(compiled.Steps()) {
      break
    }

    if _, ok := ptr[lastKey]; ok {
      ptr = ptr[lastKey].(map[string]interface{})
    } else {
      val = make(map[string]interface{})
      ptr[lastKey] = val
      ptr = val
    }
  }

  ptr[lastKey] = target

  return 0
}

func GetValue(L *lua.LState) int {
  // luaで積んだキーを取得
  key := L.ToString(1)

  // キーをjsonpathコンパイル
  compiled, err := jsonpath.Compile(key)

  if err != nil {
    panic(err)
  }

  var ptr map[string]interface{} = listHash
  var lastKey string
  for i, v := range(compiled.Steps()) {
    if v.Op() != "key" {
        continue
    }

    lastKey = v.Key()

    if i + 1 == len(compiled.Steps()) {
      break
    }

    if _, ok := ptr[lastKey]; ok {
      ptr = ptr[lastKey].(map[string]interface{})
    } else {
      L.Push(lua.LNil)
      return 1
    }
  }

  target := ptr[lastKey]

  if target == nil {
    L.Push(lua.LNil)
  } else {
    L.Push(lua.LString(target.(string)))
  }

  return 1
}

func printHashRecurse(hash map[string]interface{}) {
  for k, v := range(hash) {
    switch v.(type) {
      case map[string]interface{}:
        fmt.Print(k + ": {")
        printHashRecurse(v.(map[string]interface{}))
        fmt.Print("}")
      case string, int:
        fmt.Print(k + ": " + v.(string))
      case nil:
        fmt.Print(k + ": null")
    }
  }
}

func exist(path string) bool {
  path, err := filepath.Abs(path)
  _, err = os.Stat(path)
  return err == nil
}

func main() {
  /*****************************/
  /* luaの実行環境を作っておく */
  /*****************************/
  L := lua.NewState()

  /************************/
  /* ちゃんとクローズする */
  /************************/
  defer L.Close()

  /*****************************/
  /* luaで使用する関数をセット */
  /*****************************/
  L.SetGlobal("setValue", L.NewFunction(SetValue))
  L.SetGlobal("getValue", L.NewFunction(GetValue))
  //L.SetGlobal("hasKey",   L.NewFunction(HasKey))

  /*****************************/
  /* luaで使用する変数をセット */
  /*****************************/
  ex, err := os.Executable()
  if err != nil {
    panic(err)
  }

  exPath := filepath.Dir(ex)
  L.SetGlobal("ROOT", lua.LString(exPath))

  /*******************************/
  /* luaのスクリプトを探して実行 */
  /*******************************/
  var lua_dir string

  flag.StringVar(&lua_dir, "luadir", "./listup", "dir path for run lua script.")
  flag.Parse()

  if exist(lua_dir + "/init.lua") {
    fmt.Println("2")
    if err := L.DoFile(lua_dir + "/init.lua"); err != nil {
      panic(err)
    }
  }

  walk := func(path string, info os.FileInfo, err error) error {
    if filepath.Ext(path) == ".lua" {
      if err := L.DoFile(path); err != nil {
        panic(err)
      }
    }

    return nil
  }

  err = filepath.Walk(lua_dir, walk)

  /************/
  /* json整形 */
  /************/
  json, err := json.MarshalIndent(listHash, "", "  ")

  if err != nil {
    panic(err)
  }

  /********/
  /* 出力 */
  /********/
  fmt.Println(string(json))
}


