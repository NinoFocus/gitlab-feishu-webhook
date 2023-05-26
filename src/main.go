package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/webhooks/v6/gitlab"
	"github.com/joho/godotenv"
	"github.com/ninofocus/gitlab-feishu-webhook/src/handler"
	"github.com/ninofocus/gitlab-feishu-webhook/src/utils"
	"log"
)

func main() {
	env := flag.String("env", ".env", "Path to env file")
	flag.Parse()

	loadAndCheckEnv(env)

	server := gin.Default()

	server.GET("/ping", handlePing)
	server.POST("/gitlab/webhook", handleGitlabWebhook)

	log.Fatal(server.Run(":8083"))
}

func loadAndCheckEnv(env *string) {
	_ = godotenv.Load(*env)
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
	hook, _ := gitlab.New()

	payload, err := hook.Parse(ctx.Request, gitlab.PushEvents, gitlab.MergeRequestEvents)
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
