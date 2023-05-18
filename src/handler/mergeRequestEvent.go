package handler

import (
	"fmt"
	"github.com/go-lark/lark"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
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
	).Blue().Title(p.renderTitle()).Link(b.URL().Href(p.ObjectAttributes.URL))

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
	t := template.Must(template.New("pushEvent").Parse(`
**{{.Title}}**
{{if ne (len .Description) 0}}

{{.Description}}
{{end}}

Source Branch: {{.SourceBranch}}
Target Branch: {{.TargetBranch}}

Merge Status: {{.MergeStatus}}

Assignee: {{.Assignee.Name}}
`))
	var buf strings.Builder

	if err := t.Execute(&buf, payload.ObjectAttributes); err != nil {
		return ""
	}

	return buf.String()
}

func (payload *MyMergeRequestEventPayload) renderFooter() string {
	return payload.Repository.Name
}
