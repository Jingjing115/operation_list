package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	DataReceiveConf   DataReceiveConf `yaml:"data_receive"`
	MsgDistributeConf MsgDistributeConf `yaml:"msg_distribute"`
}

func NewConfig(configFile string) *Config {

	config := defaultConfig()

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("Could not load config because: %v\n", err)
	} else {
		if err = yaml.Unmarshal(data, &config); err != nil {
			log.Printf("Could not unmarshal config because: %v", err)
		}
	}
	return config
}

func (c *Config) String() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func defaultConfig() *Config {
	return &Config{

	}
}

type DataReceiveConf struct {
	Producer string `yaml:"producer"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Rch      string `yaml:"rch"`
	Pch      string `yaml:"pch"`
}

type MsgDistributeConf struct {
	Consumer  string `yaml:"consumer"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Database  string `yaml:"database"`
	WhiteList string `yaml:"white_list"`
}