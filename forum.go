package telebot

import (
	"encoding/json"
	"fmt"
)

type ForumTopic struct {
	ID                int64  `json:"message_thread_id"`    // Unique identifier of the forum topic
	Name              string `json:"name"`                 // Name of the topic
	IconColor         int64  `json:"icon_color"`           // Color of the topic icon in RGB format
	IconCustomEmojiId string `json:"icon_custom_emoji_id"` // Optional. Unique identifier of the custom emoji shown as the topic icon
}

func (b *Bot) CreateForumTopic(id int64, name string) (*ForumTopic, error) {
	ret, err := b.Raw("createForumTopic", struct {
		ID   int64  `json:"chat_id"`
		Name string `json:"name"`
	}{id, name})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result      *ForumTopic `json:"result"`
		Description string      `json:"description"`
	}
	if err := json.Unmarshal(ret, &resp); err != nil {
		return nil, err
	}
	if resp.Result == nil {
		return nil, fmt.Errorf("while creating topic at chat %d, named %s, error occured %s", id, name, resp.Description)
	}

	return resp.Result, nil
}
