package handler

import (
	"fmt"
	"github.com/go-lark/lark"
	"github.com/go-playground/webhooks/v6/gitlab"
	"github.com/ninofocus/gitlab-feishu-webhook/src/utils"
	"log"
	"os"
	"strings"
	"text/template"
)

type MyPushEventPayload gitlab.PushEventPayload

func HandlePushEvent(payload gitlab.PushEventPayload) {
	p := MyPushEventPayload(payload)

	bot := lark.NewNotificationBot(os.Getenv("FEISHU_BOT_WEBHOOK_URL"))

	b := lark.NewCardBuilder()

	card := b.Card(
		b.Markdown(p.renderBody()),
		b.Note().AddText(b.Text(p.renderFooter()).LarkMd()),
	).Blue().Title(p.renderTitle())

	msg := lark.NewMsgBuffer(lark.MsgInteractive)

	om := msg.Card(card.String()).Build()
	_, err := bot.PostNotificationV2(om)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Send message to feishu success")
}

func (payload *MyPushEventPayload) renderTitle() string {
	userName := payload.UserName
	totalCommitCount := payload.TotalCommitsCount

	if totalCommitCount > 1 {
		return fmt.Sprintf("%s push %d commits", userName, totalCommitCount)
	}
	return fmt.Sprintf("%s push a commit", userName)
}

func (payload *MyPushEventPayload) renderBody() string {
	t := template.Must(template.New("pushEvent").Funcs(template.FuncMap{
		"shortId": utils.GetShortCommitId,
	}).Parse(`
{{range .Commits}}
{{$shortId := shortId .ID}}
***commit [{{$shortId}}]({{.URL}})***
Author: {{.Author.Name}} ({{.Author.Email}})
Date: {{.Timestamp}}

{{.Message}}

{{end}}
`))
	var buf strings.Builder

	if err := t.Execute(&buf, payload); err != nil {
		return ""
	}

	return buf.String()
}

func (payload *MyPushEventPayload) renderFooter() string {
	branch := utils.GetBranchNameFromRef(payload.Ref)

	return fmt.Sprintf("%s > %s", payload.Repository.Name, branch)
}
