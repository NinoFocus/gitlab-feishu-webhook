package utils

import (
	"os"
	"strings"
)

func GetFeiShuBotWebhookURLFromEnv() string {
	return os.Getenv("FEISHU_BOT_WEBHOOK_URL")
}

func GetBranchNameFromRef(ref string) string {
	return strings.Replace(ref, "refs/heads/", "", 1)
}

func GetShortCommitId(commitId string) string {
	return commitId[0:7]
}
