package main

import (
	"github.com/joho/godotenv"
	"github.com/ninofocus/gitlab-feishu-webhook/src/handler"
	"github.com/ninofocus/gitlab-feishu-webhook/src/utils"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
	"log"
	"net/http"
)

func main() {
	checkEnv()

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
			go handler.HandlePushEvent(payload.(gitlab.PushEventPayload))
		case gitlab.MergeRequestEventPayload:
			go handler.HandleMergeRequestEvent(payload.(gitlab.MergeRequestEventPayload))
		}
	})
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func checkEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	webhookUrl := utils.GetFeiShuBotWebhookURLFromEnv()
	if len(webhookUrl) == 0 {
		log.Fatal("FEISHU_BOT_WEBHOOK_URL not found in env")
	}
}
