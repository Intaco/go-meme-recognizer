package watcher

import (
	"log"
	"os"
	"path/filepath"
)

type Watcher interface {
	getPhotosPaths(maxCount int)
}

// in process
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
	rootdir string
}

func NewSystemWatcher(rootdir string) (w *SystemWatcher) {
	return &SystemWatcher{rootdir}
}

func (w *SystemWatcher) getPhotosPaths(count int) (pathToPhotosList [][]string) {
	rootdir, err := os.Open(w.rootdir)
	if err != nil {
		log.Println(err)
		return
	}
	filesInfo, err := rootdir.Readdir(-1)
	if err != nil {
		log.Println(err)
		return
	}
	for _, info := range filesInfo {
		if !info.IsDir() {
			log.Println("Unexpected file, only dirs are expected in rootdir ")
			continue
		}
		var photoPaths []string
		dirPath := filepath.Join(w.rootdir, info.Name())
		dir, err := os.Open(dirPath)
		if err != nil {
			log.Println(err)
			continue
		}
		files, err := dir.Readdir(-1)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, f := range files {
			if f.IsDir() {
				log.Println("Unexpected dir, only files are expected in subdirs of rootdir")
				continue
			}
			fileName := f.Name()

			filePath := filepath.Join(dirPath, fileName())
			photoPaths = append(photoPaths, filePath)
		}
		pathToPhotosList = append(pathToPhotosList, photoPaths)
	}
	log.Println(pathToPhotosList)
	return
}
