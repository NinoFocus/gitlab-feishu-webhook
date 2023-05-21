package handler

import (
	"fmt"
	"github.com/go-lark/lark"
	"github.com/go-playground/webhooks/v6/gitlab"
	"log"
	"os"
	"strings"
	"text/template"
)

type MyMergeRequestEventPayload gitlab.MergeRequestEventPayload

func HandleMergeRequestEvent(payload gitlab.MergeRequestEventPayload) {
	p := MyMergeRequestEventPayload(payload)

	bot := lark.NewNotificationBot(os.Getenv("FEISHU_BOT_WEBHOOK_URL"))
	bot.GetTenantAccessTokenInternal(true)

	b := lark.NewCardBuilder()

	card := b.Card(
		b.Markdown(p.renderBody()),
		b.Note().AddText(b.Text(p.renderFooter()).LarkMd()),
	).Title(p.renderTitle()).Link(b.URL().Href(p.ObjectAttributes.URL))

	switch p.ObjectAttributes.Action {
	case "merge":
		card.Blue()
	case "close":
		card.Red()
	default:
		card.Blue()
	}

	msg := lark.NewMsgBuffer(lark.MsgInteractive)

	om := msg.Card(card.String()).Build()
	_, err := bot.PostNotificationV2(om)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Send message to feishu success")
}

func (payload *MyMergeRequestEventPayload) renderTitle() string {
	id := payload.ObjectAttributes.ID
	userName := payload.User.UserName
	action := payload.ObjectAttributes.Action
	return fmt.Sprintf("%s %s a merge request #%d", userName, action, id)
}

func (payload *MyMergeRequestEventPayload) renderBody() string {
	t := template.Must(template.New("mergeRequestEvent").Parse(`
**{{.ObjectAttributes.Title}}**
{{if ne (len .ObjectAttributes.Description) 0}}

{{.ObjectAttributes.Description}}
{{end}}

Source Branch: {{.ObjectAttributes.SourceBranch}}
Target Branch: {{.ObjectAttributes.TargetBranch}}

State: {{.ObjectAttributes.State}}

{{if len .Assignees}}
Assignee: {{range .Assignees}}{{.Name }} {{end}}
{{end}}
`))
	var buf strings.Builder

	if err := t.Execute(&buf, payload); err != nil {
		log.Println(err)
		return ""
	}

	return buf.String()
}

func (payload *MyMergeRequestEventPayload) renderFooter() string {
	return payload.Repository.Name
}
