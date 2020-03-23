package config

import (
	"fmt"
	"runtime"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/tkanos/gonfig"

	_ "net/http/pprof"
)

type Config struct {
	ListenPort int    `env:"TSRV_LISTENPORT"`
	TestFlag   bool   `env:"TSRV_TESTFLAG"`
	FileMap    string `env:"TSRV_FILEMAP"`
	Styles     map[string]map[string]Style
}

type Env struct {
	Config   *Config
	loglevel log.Level
}

func NewEnv(path string) *Env {
	var cfg Config
	err := gonfig.GetConf(path+"/"+"conf.json", &cfg)
	if err != nil {
		err = gonfig.GetConf(path+"/"+"conf.tpl.json", &cfg)
		checkErr(err)
	}

	loglevel := log.WarnLevel
	if cfg.TestFlag {
		loglevel = log.DebugLevel
	}

	log.SetLevel(loglevel)
	log.SetFormatter(&log.JSONFormatter{})

	return &Env{
		Config:   &cfg,
		loglevel: loglevel,
	}
}

func checkErr(err error) {
	if err != nil {
		_, filename, lineno, ok := runtime.Caller(1)
		message := ""
		if ok {
			message = fmt.Sprintf("%v:%v: %v\n", filename, lineno, err)
		}
		log.Panic(message, err)
	}
}
