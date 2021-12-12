package config

import (
	"fmt"

	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v2"
)

var ()

type RunType string

var (
	version = "dev"

	commit = "none"
	date   = "unknown"

	ProjectName string = "musicCloud-bot"
	BotToken    string
	Socks5      string

	UserName            string
	PasswdMd5           string
	PushMusicServerHost string

	// TelegramEndpoint telegram bot 服务器地址，默认为空
	TelegramEndpoint string = tb.DefaultApiURL

	// AllowUsers 允许使用bot的用户
	AllowUsers []int64
)

const (
	logo = `
	
	`
)

func AppVersionInfo() string {
	return fmt.Sprintf("version %v, commit %v, built at %v", version, commit, date)
}

// GetString get string config value by key
func GetString(key string) string {
	var value string
	if viper.IsSet(key) {
		value = viper.GetString(key)
	}

	return value
}
