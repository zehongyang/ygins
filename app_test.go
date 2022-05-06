package ygins

import (
	"github.com/zehongyang/ygins/config"
	"net/url"
	"strings"
	"testing"
)

//type Tt struct {
//	Test string `yaml:"test"`
//}
//
//func TestConfig(t *testing.T)  {
//	var st Tt
//	err := config.Load(&st)
//	t.Log(err)
//	t.Log(st.Test)
//}
//
//func TestLog(t *testing.T)  {
//	Register(handlers.Login)
//	h := Get("handlers.Login",url.Values{})
//	t.Log(h)
//}

func TestUrl(t *testing.T)  {
	var str = "wefwefewf?q=dotnet&t=111"
	query, err := url.Parse(str)
	t.Log(err)
	t.Log(query.Query())
}

type Obj struct {
	Name string
	Age int
}

func TestReflect(t *testing.T)  {

	var v = url.Values{}
	v.Set("Name","张三")
	v.Set("Age","13")
	var o Obj
	err := LoadTagStruct(&o, v)
	t.Log(err)
	t.Log(o)
}

func TestRep(t *testing.T)  {
	var str = "${ADDR}${PORT}"
	ss := strings.Split(str, "${")
	var sb strings.Builder
	for _, s := range ss {
		idx := strings.Index(s, "}")
		t.Log(s,idx)
		if idx > 0 {
			t.Log(s[:idx])
			if idx != len(s) - 1 {
				sb.WriteString(s[:idx])
				sb.WriteString(s[idx+1:])
			}else{
				sb.WriteString(strings.TrimRight(s,"}"))
			}
		}else{
			sb.WriteString(s)
		}
	}
	t.Log(sb.String())
}

func TestEnv(t *testing.T)  {
	config.IsDebug()
}