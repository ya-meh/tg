package telebot

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestBot_CreateForumTopic(t *testing.T) {
	chat, err := strconv.ParseInt(os.Getenv("TESTFORUM"), 10, 64)
	require.NoError(t, err)

	tmp, err := os.CreateTemp("", "forum_test_*.jpg")
	require.NoError(t, err)
	defer tmp.Close()

	resp, err := http.Get("https://i.ibb.co/4mKR0cv/photo-2023-09-21-09-51-53.jpg")
	require.NoError(t, err)
	defer resp.Body.Close()

	_, err = io.Copy(tmp, resp.Body)
	require.NoError(t, err)

	bot, err := NewBot(Settings{Token: os.Getenv("TESTBOT")})
	require.NoError(t, err)

	topic, err := bot.CreateForumTopic(chat, "test_topic")
	require.NoError(t, err)

	msg, err := bot.Send(&Chat{ID: chat, Topic: topic}, "test")
	require.NoError(t, err)

	_ = tmp.Close()
	msg, err = bot.Send(&Chat{ID: chat, Topic: topic}, &Photo{Caption: "test photo", File: FromDisk(tmp.Name())})
	require.NoError(t, err)

	buf, _ := json.MarshalIndent(msg, "", "  ")
	println(string(buf))

}
