package ygins

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zehongyang/ygins/config"
	"github.com/zehongyang/ygins/logger"
	"go.uber.org/zap"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	yGinMp = AppYGin{handles: make(map[string]Handler)}
	once sync.Once
	server Servers
	exitChan chan bool
	PtrErr = errors.New("is not ptr kind")
	StructErr = errors.New("is not struct kind")
)

type Handler func(values ...url.Values)gin.HandlerFunc

type Servers struct {
	Server []AppServer
}

type AppServer struct {
	Name string
	Addr string
	Middlewares []string
	RouterGroups []RouterGroup
}

type RouterGroup struct {
	Group string
	Middlewares []string
	Routers []Router
}

type Router struct {
	Path string
	Methods []string
	Handlers []string
}


type AppYGin struct {
	handles map[string]Handler
}

func Register(handlers ...Handler)  {
	for _, handler := range handlers {
		rv := reflect.ValueOf(handler)
		fn := runtime.FuncForPC(rv.Pointer())
		name := filepath.ToSlash(fn.Name())
		nameSplits := strings.Split(name, "/")
		if _,ok := yGinMp.handles[nameSplits[len(nameSplits)-1]];!ok{
			yGinMp.handles[nameSplits[len(nameSplits)-1]] = handler
		}else{
			logger.Warn("Register handler",zap.Any("hanler exist",nameSplits[len(nameSplits)-1]))
		}
	}
}

func Get(name string,v ...url.Values) gin.HandlerFunc {
	h,ok := yGinMp.handles[name]
	if !ok {
		logger.Warn("Get Handler",zap.Any("no handler named",name))
		return nil
	}
	return h(v...)
}

func Run()  {
	once.Do(func() {
		err := config.Load(&server)
		if err != nil {
			logger.Fatal("Run server",zap.Error(err),zap.Any("server",server))
		}
		wg := sync.WaitGroup{}
		for _, appServer := range server.Server {
			wg.Add(1)
			var ts = appServer
			go func() {
				defer wg.Done()
				engine := gin.New()
				if len(ts.Middlewares) > 0 {
					ts.Middlewares = distinctStringSlice(ts.Middlewares)
					for _, middleware := range ts.Middlewares {
						h := Get(middleware)
						if h != nil {
							engine.Use(h)
							logger.Info("register middleware",zap.Any("name",middleware))
						}
					}
				}
				for _, group := range ts.RouterGroups {
					routerGroup := engine.Group(group.Group)
					group.Middlewares = distinctStringSlice(group.Middlewares)
					for _, middleware := range group.Middlewares {
						h := Get(middleware)
						if h != nil {
							routerGroup.Use(h)
							logger.Info("register group middleware",zap.Any("name",middleware))
						}
					}
					for _, router := range group.Routers {
						router.Methods = distinctStringSlice(router.Methods)
						router.Handlers = distinctStringSlice(router.Handlers)
						for _, method := range router.Methods {
							var hs []gin.HandlerFunc
							for _, handler := range router.Handlers {
								h := Get(handler)
								if h != nil {
									hs = append(hs,h)
								}
							}
							if len(hs) > 0 {
								routerGroup.Handle(method,router.Path,hs...)
							}
						}
					}
				}
				err2 := engine.Run(ts.Addr)
				if err2 != nil {
					logger.Fatal("Run server",zap.Error(err),zap.Any("name",ts.Name))
				}
			}()
		}
		wg.Wait()
	})
	select {
	case <-exitChan:
		//todo exit
	}
}

func distinctStringSlice(ss []string) []string {
	var mp = make(map[string]string)
	for _, s := range ss {
		mp[s] = s
	}
	var rs []string
	for _, s := range mp {
		rs = append(rs,s)
	}
	return rs
}

func LoadTagStruct(o interface{},v url.Values,tag ...string) error {
	rv := reflect.ValueOf(o)
	if rv.Kind() != reflect.Ptr {
		return PtrErr
	}
	elem := rv.Elem()
	if elem.Kind() != reflect.Struct {
		return StructErr
	}
	rt := elem.Type()
	num := rt.NumField()
	var tg string
	if len(tag) > 0 {
		tg = tag[0]
	}
	for i := 0; i < num; i++ {
		ft := rt.Field(i)
		sv := v.Get(ft.Name)
		if len(tg) > 0 && len(sv) < 1 {
			sv = v.Get(ft.Tag.Get(tg))
		}
		field := elem.Field(i)
		if field.CanSet() && len(sv) > 0 {
			switch field.Kind() {
			case reflect.Bool:
				field.SetBool(sv == "true")
			case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
				parseInt, err := strconv.ParseInt(sv, 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(parseInt)
			case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
				parseUint, err := strconv.ParseUint(sv, 10, 64)
				if err != nil {
					return err
				}
				field.SetUint(parseUint)
			case reflect.String:
				field.SetString(sv)
			}
		}
	}
	return nil
}