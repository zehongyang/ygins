package config

import (
	"flag"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

type EnvConfig struct {
	Debug bool
	Env string
}

var (
	gConfig EnvConfig
	gData []byte
)

func init()  {
	testing.Init()
	flag.StringVar(&gConfig.Env,"env","local.env","the environment")
	flag.BoolVar(&gConfig.Debug,"d",false,"mode")
	flag.Parse()
	parseEnv()
}

func IsDebug() bool {
	return gConfig.Debug
}

func parseEnv()  {
	if len(gConfig.Env) < 1 {
		return
	}
	err := godotenv.Load(gConfig.Env)
	if err != nil {
		log.Fatalf("parseEnv err:%v",err)
	}
	loadFiles()
}

func loadFiles() {
	resource := os.Getenv("FileResources")
	if len(resource) < 1 {
		resource = "./application.yml"
	}
	resources := strings.Split(resource, ";")
	for i, s := range resources {
		data, err := ioutil.ReadFile(s)
		if err != nil {
			log.Fatalf("loadFiles err:%v",err)
		}
		if i > 0 {
			gData = append(gData,'\n')
		}
		gData = append(gData,data...)
	}
	replaceEnvVariable()
}

func replaceEnvVariable()  {
	str := string(gData)
	if len(str) < 1 {
		return
	}
	ss := strings.Split(str, "${")
	var sb strings.Builder
	for _, s := range ss {
		idx := strings.Index(s, "}")
		if idx > 0 {
			var k = s[:idx]
			envK := os.Getenv(k)
			if len(envK) > 0 {
				sb.WriteString(envK)
			}else{
				sb.WriteString(k)
			}
			if idx != len(s) - 1 {
				sb.WriteString(s[idx+1:])
			}
		}else{
			sb.WriteString(s)
		}
	}
	gData = []byte(sb.String())
}

func Load(o interface{}) error {
	return yaml.Unmarshal(gData, o)
}
