package bot

import (
	"musicCloud-bot/config"
	"musicCloud-bot/log"

	tb "gopkg.in/tucnak/telebot.v2"
)

//SendError send error user
func SendError(c *tb.Chat) {
	_, _ = B.Send(c, "请输入正确的指令！")
}

// CheckAdmin check user is admin of group/channel
func CheckAdmin(upd *tb.Update) bool {

	if upd.Message != nil {
		if HasAdminType(upd.Message.Chat.Type) {
			adminList, _ := B.AdminsOf(upd.Message.Chat)
			for _, admin := range adminList {
				if admin.User.ID == upd.Message.Sender.ID {
					return true
				}
			}

			return false
		}

		return true
	} else if upd.Callback != nil {
		if HasAdminType(upd.Callback.Message.Chat.Type) {
			adminList, _ := B.AdminsOf(upd.Callback.Message.Chat)
			for _, admin := range adminList {
				if admin.User.ID == upd.Callback.Sender.ID {
					return true
				}
			}

			return false
		}

		return true
	}
	return false
}

// IsUserAllowed check user is allowed to use bot
func isUserAllowed(upd *tb.Update) bool {
	if upd == nil {
		return false
	}

	var userID int64

	if upd.Message != nil {
		userID = int64(upd.Message.Sender.ID)
	} else if upd.Callback != nil {
		userID = int64(upd.Callback.Sender.ID)
	} else {
		return false
	}

	if len(config.AllowUsers) == 0 {
		return true
	}

	for _, allowUserID := range config.AllowUsers {
		if allowUserID == userID {
			return true
		}
	}

	log.Infow("user not allowed", "userID", userID)
	return false
}

func userIsAdminOfGroup(userID int, groupChat *tb.Chat) (isAdmin bool) {

	adminList, err := B.AdminsOf(groupChat)
	isAdmin = false

	if err != nil {
		return
	}

	for _, admin := range adminList {
		if userID == admin.User.ID {
			isAdmin = true
		}
	}
	return
}

// UserIsAdminChannel check if the user is the administrator of channel
func UserIsAdminChannel(userID int, channelChat *tb.Chat) (isAdmin bool) {
	adminList, err := B.AdminsOf(channelChat)
	isAdmin = false

	if err != nil {
		return
	}

	for _, admin := range adminList {
		if userID == admin.User.ID {
			isAdmin = true
		}
	}
	return
}

// HasAdminType check if the message is sent in the group/channel environment
func HasAdminType(t tb.ChatType) bool {
	hasAdmin := []tb.ChatType{tb.ChatGroup, tb.ChatSuperGroup, tb.ChatChannel, tb.ChatChannelPrivate}
	for _, n := range hasAdmin {
		if t == n {
			return true
		}
	}
	return false
}

// GetMentionFromMessage get message mention
func GetMentionFromMessage(m *tb.Message) (mention string) {
	if m.Text != "" {
		for _, entity := range m.Entities {
			if entity.Type == tb.EntityMention {
				if mention == "" {
					mention = m.Text[entity.Offset : entity.Offset+entity.Length]
					return
				}
			}
		}
	} else {
		for _, entity := range m.CaptionEntities {
			if entity.Type == tb.EntityMention {
				if mention == "" {
					mention = m.Caption[entity.Offset : entity.Offset+entity.Length]
					return
				}
			}
		}
	}
	return
}

// GetURLAndMentionFromMessage get URL and mention from message
func GetURLAndMentionFromMessage(m *tb.Message) (url string, mention string) {
	for _, entity := range m.Entities {
		if entity.Type == tb.EntityMention {
			if mention == "" {
				mention = m.Text[entity.Offset : entity.Offset+entity.Length]

			}
		}

		if entity.Type == tb.EntityURL {
			if url == "" {
				url = m.Text[entity.Offset : entity.Offset+entity.Length]
			}
		}
	}

	return
}
