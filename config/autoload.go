package config

import (
	"flag"
	"fmt"
	"musicCloud-bot/pushmusic"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

func init() {

	workDirFlag := flag.String("d", "./", "work directory of musicCloud-bot")
	configFile := flag.String("c", "", "config file of musicCloud-bot")
	printVersionFlag := flag.Bool("v", false, "prints musicCloud-bot version")

	flag.Parse()

	if *printVersionFlag {
		// print version
		fmt.Printf(AppVersionInfo())
		os.Exit(0)
	}

	workDir := filepath.Clean(*workDirFlag)

	if *configFile != "" {
		viper.SetConfigFile(*configFile)
	} else {
		viper.SetConfigFile(filepath.Join(workDir, "config.yml"))
	}

	fmt.Println(logo)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	BotToken = viper.GetString("bot_token")
	Socks5 = viper.GetString("socks5")

	UserName = viper.GetString("netEasyMusicServer.userName")
	PasswdMd5 = viper.GetString("netEasyMusicServer.passwdMd5")
	PushMusicServerHost = viper.GetString("netEasyMusicServer.pushMusicServerHost")

	ret, err := pushmusic.Login(PushMusicServerHost, UserName, PasswdMd5)
	if !ret {
		panic(fmt.Errorf("Login the music server fail: %s", err))
	}

	if viper.IsSet("allowed_users") {
		intAllowUsers := viper.GetStringSlice("allowed_users")
		for _, useIDStr := range intAllowUsers {
			userID, err := strconv.ParseInt(useIDStr, 10, 64)
			if err != nil {
				panic(fmt.Errorf("Fatal error config file: %s", err))
			}
			AllowUsers = append(AllowUsers, userID)
		}
	}

	if viper.IsSet("telegram.endpoint") {
		TelegramEndpoint = viper.GetString("telegram.endpoint")
	}

}

func getInt(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}
