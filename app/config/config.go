package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var gManager *Manager

// Manager manage configs.
type Manager struct {
	serve *Serve
}

// GetServe get serve config.
func GetServe() Serve {
	return *gManager.serve
}

// LoadConfig load config.
func LoadConfig(path string) error {
	// 读取基本配置
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	serve := Serve{}
	err = yaml.Unmarshal(data, &serve)
	if err != nil {
		return err
	}

	// 设置全局配置
	gManager = &Manager{
		serve: &serve,
	}
	return nil
}
