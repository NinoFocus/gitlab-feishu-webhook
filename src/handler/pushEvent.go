package handler

import (
	"fmt"
	"github.com/go-lark/lark"
	"github.com/ninofocus/gitlab-feishu-webhook/src/utils"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
	"log"
	"os"
	"strings"
	"text/template"
)

func HandlePushEvent(pushEvent gitlab.PushEventPayload) {
	bot := lark.NewNotificationBot(os.Getenv("FEISHU_BOT_WEBHOOK_URL"))
	bot.GetTenantAccessTokenInternal(true)

	b := lark.NewCardBuilder()

	card := b.Card(
		b.Markdown(renderBody(pushEvent)),
		b.Note().AddText(b.Text(renderFooter(pushEvent)).LarkMd()),
	).Blue().Title(renderTitle(pushEvent))

	msg := lark.NewMsgBuffer(lark.MsgInteractive)

	om := msg.Card(card.String()).Build()
	_, err := bot.PostNotificationV2(om)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Send message to feishu success")
}

func renderTitle(pushEvent gitlab.PushEventPayload) string {
	userName := pushEvent.UserName
	totalCommitCount := pushEvent.TotalCommitsCount

	if totalCommitCount > 1 {
		return fmt.Sprintf("%s push %d commits", userName, totalCommitCount)
	}
	return fmt.Sprintf("%s push a commit", userName)
}

func renderBody(pushEvent gitlab.PushEventPayload) string {
	t := template.Must(template.New("pushEvent").Funcs(template.FuncMap{
		"shortId": utils.GetShortCommitId,
	}).Parse(`
{{range .Commits}}
{{$shortId := shortId .ID}}
*commit [{{$shortId}}]({{.URL}})*
Author: {{.Author.Name}} ({{.Author.Email}})

{{.Message}}

{{end}}
`))
	var buf strings.Builder

	if err := t.Execute(&buf, pushEvent); err != nil {
		return ""
	}

	return buf.String()
}

func renderFooter(pushEvent gitlab.PushEventPayload) string {
	branch := utils.GetBranchNameFromRef(pushEvent.Ref)

	return fmt.Sprintf("%s > %s", pushEvent.Repository.Name, branch)
}
