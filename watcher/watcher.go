package main

import (
	"fmt"
	"log"
	"os"
)

/*
func main() {
	watcher := NewVkWatcher(VkAppId, VkAppSecretKey, -67580761)
	posts := watcher.getPhotosPaths(110)
	fmt.Println(posts)
	fmt.Println(len(posts))
}
*/

var APP_LOGGER = log.New(os.Stderr, "", 0)

type Watcher interface {
	getPhotosPaths(maxCount int)
}

type FacebookWatcher struct {
}

func (w *FacebookWatcher) getPhotosPaths(count int) (pathToPhotosList [][]string) {
	return
}

type TwitterWatcher struct {
}

func (w *TwitterWatcher) getPhotosPaths(count int) (pathToPhotosList [][]string) {
	return
}

type TelegramWatcher struct {
}

func (w *TelegramWatcher) getPhotosPaths(count int) (pathToPhotosList [][]string) {
	return
}

type SystemWatcher struct {
}

func (w *SystemWatcher) getPhotosPaths(count int) (pathToPhotosList [][]string) {
	return
}
