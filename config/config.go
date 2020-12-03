package config

import (
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

//Config project config
type Config struct {
	Mongo `yaml:"mongo"`
	Port  int    `yaml:"port"`
	Db    string `yaml:"db"`
}

//Mongo config
type Mongo struct {
	MongoURL string `yaml:"url"`
	Poolsize uint64 `yaml:"poolsize"`
	DbName   string `yaml:"dbname"`
}

var (
	//Cfg config
	Cfg    Config
	inited bool
	m      = &sync.Mutex{}
)

//Init init
func Init(configpath string) error {
	m.Lock()
	defer m.Unlock()
	if inited {
		return nil
	}
	inited = true
	f, err := os.Open(configpath)
	if err != nil {
		return err
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, &Cfg)
	if err == nil {
		//纳秒转换成秒
	}
	return err
}
