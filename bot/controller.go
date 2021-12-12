package bot

import (
	"musicCloud-bot/config"
	"musicCloud-bot/log"
	musiccore "musicCloud-bot/musiccore"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

func helpCmdCtr(m *tb.Message) {
	message := `
	本助手旨在帮你:拷贝YouTube或者哔哩哔哩的视频链接发到这里，可以帮你转换成mp3音频推送到你的网易云音乐个人云盘里。
`
	_, _ = B.Send(m.Chat, message)
}

func versionCmdCtr(m *tb.Message) {
	_, _ = B.Send(m.Chat, config.AppVersionInfo())
}

func myIDCmdCtr(m *tb.Message) {
	_, _ = B.Send(m.Chat, strconv.Itoa(m.Sender.ID))
}

func textCtr(m *tb.Message) {
	log.Debug(m.Text)
	req := musiccore.NewReqPack(m.Text, m, sendMsg)
	go req.HandleMsg()
}

func sendMsg(m *tb.Message, msg string) {
	log.Println(msg)
	_, _ = B.Send(m.Chat, msg)
}
