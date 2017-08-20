package main

import (
	"encoding/json"
)

type vkResponse struct {
	Response []json.RawMessage `json:"response"`
}

type vkPost struct {
	Attachment  vkAttachment   `json:"attachment"`
	Attachments []vkAttachment `json:"attachments"`
	Comments    vkComments     `json:"comments"`
	Date        int64          `json:"date"`
	FromID      int64          `json:"from_id"`
	Id          int64          `json:"id"`
	IsPinned    int64          `json:"is_pinned"`
	Likes       vkLikes        `json:"likes"`
	MarkedAsAds int64          `json:"marked_as_ads"`
	Media       vkMedia        `json:"media"`
	Online      int64          `json:"online"`
	PostSource  vkPostSource   `json:"post_source"`
	PostType    string         `json:"post_type"`
	ReplyCount  int64          `json:"reply_count"`
	Reposts     vkReposts      `json:"reposts"`
	Text        string         `json:"text"`
	ToID        int64          `json:"to_id"`
}

type vkPhoto struct {
	AccessKey string `json:"access_key"`
	Aid       int64  `json:"aid"`
	Created   int64  `json:"created"`
	Height    int64  `json:"height"`
	OwnerID   int64  `json:"owner_id"`
	Pid       int64  `json:"pid"`
	PostID    int64  `json:"post_id"`
	Src       string `json:"src"`
	SrcBig    string `json:"src_big"`
	SrcSmall  string `json:"src_small"`
	SrcXbig   string `json:"src_xbig"`
	SrcXxbig  string `json:"src_xxbig"`
	SrcXxxbig string `json:"src_xxxbig"`
	Text      string `json:"text"`
	UserID    int64  `json:"user_id"`
	Width     int64  `json:"width"`
}

type vkLikes struct {
	CanLike    int64 `json:"can_like"`
	CanPublish int64 `json:"can_publish"`
	Count      int64 `json:"count"`
	UserLikes  int64 `json:"user_likes"`
}

type vkComments struct {
	CanPost int64 `json:"can_post"`
	Count   int64 `json:"count"`
}

type vkReposts struct {
	Count        int64 `json:"count"`
	UserReposted int64 `json:"user_reposted"`
}

type vkMedia struct {
	ItemID   int64  `json:"item_id"`
	OwnerID  int64  `json:"owner_id"`
	ThumbSrc string `json:"thumb_src"`
	Kind     string `json:"type"`
}

type vkAttachment struct {
	Photo vkPhoto `json:"photo"`
	Kind  string  `json:"type"`
}

type vkPostSource struct {
	Kind string `json:"type"`
}
