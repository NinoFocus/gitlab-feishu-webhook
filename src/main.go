package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ninofocus/gitlab-feishu-webhook/src/handler"
	"github.com/ninofocus/gitlab-feishu-webhook/src/utils"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
	"log"
)

func main() {
	godotenv.Load()
	checkEnv()

	server := gin.Default()

	server.GET("/ping", handlePing)
	server.POST("/gitlab/webhook", handleGitlabWebhook)

	server.Run(":8083")
}

func checkEnv() {
	webhookUrl := utils.GetFeiShuBotWebhookURLFromEnv()
	if len(webhookUrl) == 0 {
		log.Fatal("FEISHU_BOT_WEBHOOK_URL not found in env")
	}
}

func handlePing(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

func handleGitlabWebhook(ctx *gin.Context) {
	r := ctx.Request

	event := r.Header.Get("X-Gitlab-Event")
	if len(event) == 0 {
		log.Println("Empty X-Gitlab-Event")
		return
	}

	hook, _ := gitlab.New()

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
}
