package bot

import (
	"net/http"
	"time"

	"musicCloud-bot/config"
	"musicCloud-bot/log"

	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	httpClient *http.Client
	// B telebot
	B *tb.Bot
)

func init() {
	poller := &tb.LongPoller{Timeout: 10 * time.Second}
	spamProtected := tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {
		if !isUserAllowed(upd) {
			// 检查用户是否可以使用bot
			return false
		}

		if !CheckAdmin(upd) {
			return false
		}
		return true
	})
	log.Infow("init telegram bot",
		"token", config.BotToken,
		"endpoint", config.TelegramEndpoint,
	)
	//http init
	httpTransport := &http.Transport{}
	httpClient = &http.Client{Transport: httpTransport, Timeout: 15 * time.Second}
	// set proxy
	if config.Socks5 != "" {
		log.Infow("enable proxy",
			"socks5", config.Socks5,
		)

		dialer, err := proxy.SOCKS5("tcp", config.Socks5, nil, proxy.Direct)
		if err != nil {
			log.Fatal("Error creating dialer, aborting.")
		}
		httpTransport.Dial = dialer.Dial
	}
	// create bot
	var err error
	B, err = tb.NewBot(tb.Settings{
		URL:    config.TelegramEndpoint,
		Token:  config.BotToken,
		Poller: spamProtected,
		Client: httpClient,
	})

	if err != nil {
		log.Fatal(err)
		return
	}
}

//Start bot
func Start() {
	setCommands()
	setHandle()
	B.Start()
}

func setCommands() {
	// 设置bot命令提示信息
	commands := []tb.Command{
		{"help", "使用帮助"},
		{"version", "bot版本"},
		{"myid", "你的聊天ID"},
	}

	if err := B.SetCommands(commands); err != nil {
		log.Errorw("set bot commands failed", "error", err.Error())
	}
}

func setHandle() {

	B.Handle("/help", helpCmdCtr)

	B.Handle("/version", versionCmdCtr)

	B.Handle("/myid", myIDCmdCtr)

	B.Handle(tb.OnText, textCtr)

}
