package telebot

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestBot_CreateForumTopic(t *testing.T) {
	bot, err := NewBot(Settings{Token: os.Getenv("TESTBOT")})
	require.NoError(t, err)

	topic, err := bot.CreateForumTopic(-1001914735047, "test_topic")
	require.NoError(t, err)

	buf, _ := json.MarshalIndent(topic, "", "  ")
	println(string(buf))

	msg, err := bot.Send(&Chat{ID: -1001914735047, Topic: topic}, "test")
	require.NoError(t, err)

	buf, _ = json.MarshalIndent(msg, "", "  ")
	println(string(buf))

}
