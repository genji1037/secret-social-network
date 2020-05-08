package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var gManager *Manager

// 配置管理器
type Manager struct {
	serve *Serve
}

// 获取服务配置
func GetServe() Serve {
	return *gManager.serve
}

// 加载配置文件
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
