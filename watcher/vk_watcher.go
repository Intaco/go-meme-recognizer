package watcher

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"log"
)

const VK_APP_TIMEOUT uint = 10

const VkAppId int = 6151326

const VkAppSecretKey string = "FeJG9L5Fb5WkZsPD0ShV"

var VkGroups = [...]int{-67580761}

const vkApiWallMaxAmount int = 100

type VkWatcher struct {
	app_id        int
	client_secret string
	access_token  string // should be initialized
	aim_id        int    // group or user id, we want get posts from it's wall
}

func NewVkWatcher(app_id int, client_secret string, aim_id int) *VkWatcher {
	p := new(VkWatcher)
	p.app_id = app_id
	p.client_secret = client_secret
	p.initAccessToken()
	p.aim_id = aim_id
	return p
}

func (w *VkWatcher) initAccessToken() {
	w.access_token = "0b463c42b423164cc496ea925b51ec75de2b474a47344da8c468eade696249c2d711bf5bcba36bdc84428"
	// fuck all creators and maintainers of vk, especially person(s) taking care about dev manual
}

func (w *VkWatcher) makeVkWallRequest(count int, offset int) (response []byte, err error) {
	timeout := time.Duration(VK_APP_TIMEOUT) * time.Second
	client := http.Client{Timeout: timeout}
	resp, err := client.PostForm("https://api.vk.com/method/wall.get",
		url.Values{"access_token": {w.access_token},
			"owner_id": {strconv.Itoa(w.aim_id)},
			"count":    {strconv.Itoa(count)},
			"filter":   {"owner"},
			"offset":   {strconv.Itoa(offset)}})
	if err != nil {
		log.Println(err)
		return
	}
	response, err = ioutil.ReadAll(resp.Body)
	return
}

func (w *VkWatcher) responseToPosts(response []byte) (posts []vkPost, err error) {
	var p vkResponse
	err = json.Unmarshal(response, &p)
	if err != nil {
		log.Println(err)
		return
	}
	for _, raw_post := range p.Response[1:] {
		post, err := w.rawJsonToPost(raw_post)
		if err != nil {
			log.Println(err)
			continue
		}
		posts = append(posts, post)
	}
	return
}

func (w *VkWatcher) rawJsonToPost(rawPost []byte) (post vkPost, err error) {
	err = json.Unmarshal(rawPost, &post)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (w *VkWatcher) getPostsChunk(count int, offset int) (posts []vkPost, err error) {
	response, err := w.makeVkWallRequest(count, offset)
	if err != nil {
		log.Println(err)
		return
	}
	posts, err = w.responseToPosts(response)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (w *VkWatcher) getPosts(count int) (posts []vkPost) {
	for offset := 0; offset < count; offset += 100 {
		var limit int
		if offset+vkApiWallMaxAmount < count {
			limit = vkApiWallMaxAmount
		} else {
			limit = count - offset
		}
		posts_chunk, err := w.getPostsChunk(limit, offset)
		if err != nil {
			continue
		}
		posts = append(posts, posts_chunk...)
	}
	return
}

func (w *VkWatcher) photosFromPosts(posts []vkPost) (pathToPhotosList [][]string) {
	for _, post := range posts {
		pathToPhotos := w.photoFromPost(post)
		if len(pathToPhotos) != 0 {
			pathToPhotosList = append(pathToPhotosList, pathToPhotos)
		}
	}
	return
}

func (w *VkWatcher) photoFromPost(post vkPost) (pathToPhotoList []string) {
	for _, attachment := range post.Attachments {
		if attachment.Kind != "photo" {
			continue
		}
		photo := attachment.Photo
		var pathToPhoto string
		if photo.SrcXxxbig != "" {
			pathToPhoto = photo.SrcXxxbig
		} else if photo.SrcXxbig != "" {
			pathToPhoto = photo.SrcXxbig
		} else if photo.SrcXbig != "" {
			pathToPhoto = photo.SrcXbig
		} else if photo.SrcBig != "" {
			pathToPhoto = photo.SrcBig
		} else if photo.SrcSmall != "" {
			pathToPhoto = photo.SrcSmall
		} else {
			continue
		}
		pathToPhotoList = append(pathToPhotoList, pathToPhoto)
	}
	return
}

func (w *VkWatcher) getPhotosPaths(count int) (pathToPhotosList [][]string) {
	posts := w.getPosts(count)
	pathToPhotosList = w.photosFromPosts(posts)
	return
}
