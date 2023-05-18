package handler

import (
	"github.com/go-lark/lark"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
	"os"
)

func HandleMergeRequestEvent(mergeRequestEvent gitlab.MergeRequestEventPayload) {
	bot := lark.NewNotificationBot(os.Getenv("FEISHU_BOT_WEBHOOK_URL"))
	bot.PostNotificationV2(lark.NewMsgBuffer(lark.MsgText).Text("Merge Request Event").Build())
}
