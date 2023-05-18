package main

import (
	"fmt"
	"github.com/go-lark/lark"
	"github.com/joho/godotenv"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	hook, _ := gitlab.New()

	http.HandleFunc("/gitlab/webhook", func(w http.ResponseWriter, r *http.Request) {
		event := r.Header.Get("X-Gitlab-Event")
		if len(event) == 0 {
			log.Println("Empty X-Gitlab-Event")
			return
		}

		gitLabEvent := gitlab.Event(event)
		log.Printf("Received: [%s]", gitLabEvent)

		payload, err := hook.Parse(r, gitlab.PushEvents, gitlab.MergeRequestEvents)
		if err != nil {
			log.Println(err)
			return
		}
		switch payload.(type) {
		case gitlab.PushEventPayload:
			go handlePushEvent(payload.(gitlab.PushEventPayload))
		case gitlab.MergeRequestEventPayload:
			go handleMergeRequestEvent(payload.(gitlab.MergeRequestEventPayload))
		}
	})
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func handlePushEvent(pushEvent gitlab.PushEventPayload) {
	bot := lark.NewNotificationBot(os.Getenv("FEISHU_BOT_WEBHOOK_URL"))
	bot.GetTenantAccessTokenInternal(true)

	b := lark.NewCardBuilder()

	card := b.Card(
		b.Markdown(renderBody(pushEvent)),
		b.Note().AddText(b.Text(renderFooter(pushEvent)).LarkMd()),
	).Blue().Title(renderTitle(pushEvent))

	msg := lark.NewMsgBuffer(lark.MsgInteractive)
	log.Println(card.String())
	om := msg.Card(card.String()).Build()
	response, err := bot.PostNotificationV2(om)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(response)
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
	t := template.Must(template.New("pushEvent").Parse(`
{{range .Commits}}
commit [{{.ID}}]({{.URL}}) 
Author: {{.Author.Name}} <{{.Author.Email}}>

{{.Message}}

{{end}}
`))
	var buf strings.Builder

	if err := t.Execute(&buf, pushEvent); err != nil {
		panic(err)
	}
	log.Println(buf.String())
	return buf.String()
}

func renderFooter(pushEvent gitlab.PushEventPayload) string {
	branch := strings.Replace(pushEvent.Ref, "refs/heads/", "", 1)

	return fmt.Sprintf("%s > %s", pushEvent.Repository.Name, branch)
}

func handleMergeRequestEvent(mergeRequestEvent gitlab.MergeRequestEventPayload) {
	bot := lark.NewNotificationBot(os.Getenv("FEISHU_BOT_WEBHOOK_URL"))
	bot.PostNotificationV2(lark.NewMsgBuffer(lark.MsgText).Text("Merge Request Event").Build())
}
