package config

import log "github.com/sirupsen/logrus"

// 服务配置
type Serve struct {
	Host      string       `yaml:"host"`      // 主机地址
	Port      int          `yaml:"port"`      // 端口号
	Level     Level        `yaml:"level"`     // 日志级别
	MySQL     MySQL        `yaml:"mysql"`     // MySQL数据库
	Open      OpenPlatform `yaml:"open"`      // 开放平台
	Consensus Consensus    `yaml:"consensus"` // 共识模块
}

// 日志级别
type Level string

func (level Level) Value() log.Level {
	lvl, err := log.ParseLevel(string(level))
	if err != nil {
		return log.InfoLevel
	}
	return lvl
}

// MySQL配置
type MySQL struct {
	Host         string `yaml:"host"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Database     string `yaml:"database"`
	Charset      string `yaml:"charset"`
	AlwayMigrate bool   `yaml:"alway_migrate"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type OpenPlatform struct {
	BaseURL   string `yaml:"base_url"`
	SecretKey string `yaml:"secret_key"`
}

type Consensus struct {
	Token         string `yaml:"token"`          // 建立连接使用的币种
	PaymentRemark string `yaml:"payment_remark"` // 支付备注
}
