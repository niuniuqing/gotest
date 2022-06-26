package demo

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type Error struct {
	errCode uint8
}

func (e *Error) Error() string {
	switch e.errCode {
	case 1:
		return "file not found"
	default:
		return "unknown error"
	}
}
func checkError(err error) {
	if err != nil && !reflect.ValueOf(err).IsNil() {
		panic(err)
	}
}

func TestDemo(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var e *Error
		checkError(e)
	})

	t.Run("map test", func(t *testing.T) {
		m := map[string]int{
			"小来": 80,
			"博清": 18,
		}
		i, _ := m["小来"]
		fmt.Println(i)
	})

	t.Run("test", func(t *testing.T) {
		var str = "liuboqing"
		bytes := []byte(str)
		md5Ctx := md5.New()
		i, err := md5Ctx.Write(bytes)
		fmt.Println("i", i, "err", err)
		fmt.Printf("%x\n", md5Ctx.Sum(nil))
		fmt.Printf("%x-------", md5Ctx.Sum([]byte("123")))
	})

	t.Run("json 转换", func(t *testing.T) {
		stu := &student{
			Name:  "张三",
			Age:   17,
			Email: "2183902@qq.com",
		}
		by, _ := json.Marshal(stu)
		fmt.Println("字符串：", string(by))
	})

	t.Run("转为struct", func(t *testing.T) {
		str := `{"name":"张三","age":17,"email":"2183902@qq.com"}`
		stu1 := &student{}
		json.Unmarshal([]byte(str), stu1)
		fmt.Println(stu1)
	})
}

type student struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}
