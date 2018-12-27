package wxwork

import (
	"context"
)

// Weixin work message type.
type MsgType string

const (
	MsgTypeText              MsgType = "text"
	MsgTypeImage                     = "image"
	MsgTypeVoice                     = "voice"
	MsgTypeVideo                     = "video"
	MsgTypeFile                      = "file"
	MsgTypeNews                      = "news"
	MsgTypeTextCard                  = "textcard"
	MsgTypeMarkdown                  = "markdown"
	MsgTypeMiniProgramNotice         = "miniprogram_notice"
)

// Message represents a weixin work message.
type Message interface {
	// Get type of this message.
	Type() MsgType
}

// Text message.
type Text struct {
	Content string `json:"content"`
}

func (t *Text) Type() MsgType {
	return MsgTypeText
}

// Image message.
type Image struct {
	MediaID string `json:"media_id"`
}

func (i *Image) Type() MsgType {
	return MsgTypeImage
}

// Voice message.
type Voice struct {
	MediaID string `json:"media_id"`
}

func (v *Voice) Type() MsgType {
	return MsgTypeVoice
}

// Video message.
type Video struct {
	MediaID     string `json:"media_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (v *Video) Type() MsgType {
	return MsgTypeVideo
}

// File message.
type File struct {
	MediaID string `json:"media_id"`
}

func (f *File) Type() MsgType {
	return MsgTypeFile
}

// News message.
type News struct {
	Articles []*NewsArticle `json:"articles"`
}

func (n *News) Type() MsgType {
	return MsgTypeNews
}

// Article of news message.
type NewsArticle struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PicURL      string `json:"picurl"`
}

// Text card message.
type TextCard struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BtnTxt      string `json:"btntxt,omitempty"`
}

func (tc *TextCard) Type() MsgType {
	return MsgTypeTextCard
}

// Markdown message.
type Markdown struct {
	Content string `json:"content"`
}

func (m *Markdown) Type() MsgType {
	return MsgTypeMarkdown
}

// Mini program notice message.
type MiniProgramNotice struct {
	AppID             string         `json:"appid"`
	Page              string         `json:"page"`
	Title             string         `json:"title"`
	Description       string         `json:"description"`
	EmphasisFirstItem bool           `json:"emphasis_first_item"`
	ContentItems      []KeyValuePair `json:"content_items"`
}

func (m *MiniProgramNotice) Type() MsgType {
	return MsgTypeMiniProgramNotice
}

// Key value pair.
type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Message service provides methods for send message.
type MessageService service

// Send sends a message to targets.
//
// Weixin Work API docs: https://work.weixin.qq.com/api/doc#90000/90135/90236
func (s *MessageService) Send(ctx context.Context, agentID int, targets *TargetSet, m Message, opt *SendOptions) (*SendResult, *Response, error) {
	return nil, nil, nil
}

// Options for send weixin work messages.
type SendOptions struct {
	// Safe message.
	Safe bool
}

// Result of send action.
type SendResult struct {
	InvalidTargets TargetSet
}

type rawSendRequest struct {
	ToUser  UserSet  `json:"touser,omitempty"`
	ToParty PartySet `json:"toparty,omitempty"`
	ToTag   TagSet   `json:"totag,omitempty"`

	AgentID int     `json:"agentid"`
	MsgType MsgType `json:"msgtype"`

	Text              Text              `json:"text,omitempty"`
	Image             Image             `json:"image,omitempty"`
	Voice             Voice             `json:"voice,omitempty"`
	Video             Video             `json:"video,omitempty"`
	File              File              `json:"file,omitempty"`
	News              News              `json:"news,omitempty"`
	TextCard          TextCard          `json:"textcard,omitempty"`
	Markdown          Markdown          `json:"markdown,omitempty"`
	MiniProgramNotice MiniProgramNotice `json:"miniprogram_notice,omitempty"`
	Safe              bool              `json:"safe,omitempty"`
}

type rawSendResponse struct {
	InvalidUser  UserSet  `json:"invaliduser,omitempty"`
	InvalidParty PartySet `json:"invalidparty,omitempty"`
	InvalidTag   TagSet   `json:"invalidtag,omitempty"`
}
