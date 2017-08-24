package fetcher

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Task interface that requires Execute() method
type Task interface {
	Execute()
}

// Pool pool
type Pool struct {
	mu    sync.Mutex
	size  int
	tasks chan Task
	kill  chan struct{}
	wg    sync.WaitGroup
}

// DefaultPoolTaskChannelSize default size of pool
const DefaultPoolTaskChannelSize = 128

// NewPool create new pool of given size
func NewPool(size int) *Pool {
	pool := &Pool{
		tasks: make(chan Task, DefaultPoolTaskChannelSize),
		kill:  make(chan struct{}),
	}
	pool.Resize(size)
	return pool
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok { // if chanel is closed, die
				return
			}
			task.Execute()
		case <-p.kill: // it's time to die!
			return
		}
	}
}

// Resize resize pool to specified size
func (p *Pool) Resize(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.size < n {
		p.size++
		p.wg.Add(1)
		go p.worker()
	}
	for p.size > n {
		p.size--
		p.kill <- struct{}{}
	}
}

// Close close pool
func (p *Pool) Close() {
	close(p.tasks)
}

// Wait wait for all tasks in pool
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Exec schedule Task for execution
func (p *Pool) Exec(e Task) {
	p.tasks <- e
}

// DownloadTask represents download task for fetcher
type DownloadTask struct {
	query            DownloadQuery
	downloadWaitTime int
}

// NewDownloadTask create download task
func NewDownloadTask(query DownloadQuery, downloadWaitTime int,
	verbose int) (t *DownloadTask) {
	t = &DownloadTask{query, DefaultPoolTaskChannelSize}
	return
}

func (t *DownloadTask) download() (err error) {
	log.Printf("Started downloading %s\n\tto %s\n", t.query.url, t.query.GetFilePath())

	fd, err := t.query.prepare()
	if err != nil {
		log.Println(err)
		return
	}

	timeout := time.Duration(t.downloadWaitTime) * time.Second
	client := http.Client{Timeout: timeout}

	response, err := client.Get(t.query.url)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(fd, response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Finished downloading %s\n\t%s\n", t.query.url,
		t.query.GetFilePath())
	t.query.markAsDone()
	return
}

// Execute perform DownloadTask
func (t DownloadTask) Execute() {
	t.download()
}

// DownloadQuery query downloads
type DownloadQuery struct {
	url      string
	dirPath  string
	filename string
	isDone   bool
}

// NewDownloadQuery create download query
func NewDownloadQuery(url string, dirPath string,
	filename string) DownloadQuery {

	return DownloadQuery{url, dirPath, filename, false}
}

// GetFilePath get query target filepath
func (q *DownloadQuery) GetFilePath() (filePath string) {
	filePath = filepath.Join(q.dirPath, q.filename)
	return
}

// IsDone is query done
func (q *DownloadQuery) IsDone() bool {
	return q.isDone
}

func (q *DownloadQuery) markAsDone() {
	q.isDone = true
}

func (q *DownloadQuery) prepare() (fd *os.File, err error) {
	if err = os.MkdirAll(q.dirPath, 0777); err != nil {
		log.Println(err)
		return
	}
	filePath := q.GetFilePath()
	if fd, err = os.Create(filePath); err != nil {
		log.Println(err)
		return
	}

	return
}

// Fetcher fetch download queries
type Fetcher struct {
	concurrency int
	wait_time   int
}

// NewFetcher create new fetcher
func NewFetcher(concurrency int, waitTime int) (f *Fetcher) {
	f = &Fetcher{concurrency, waitTime}
	return
}

// Download execute download queries. Start from startIndex-th query
func (f *Fetcher) Download(downloads []DownloadQuery, startIndex int) {
	p := pool.NewPool(f.concurrency)
	for _, query := range downloads {
		pool.Exec(NewDownloadTask(query,
			f.waitTime, 1))
	}
	p.Close()
	p.Wait()
}

// MakeQueryFromUrlsList make DownloadQueries from urls
func MakeQueryFromUrlsList(rootdir string, urlsList [][]string) (downloads []DownloadQuery) {
	for index1, urls := range urlsList {
		prefix := strconv.Itoa(index1)
		dirPath := filepath.Join(rootdir, prefix)
		for index2, url := range urls {
			filename := prefix + "_" + strconv.Itoa(index2) +
				url[strings.LastIndex(url, "."):]
			downloads = append(downloads, NewDownloadQuery(url, dirPath, filename))
		}
	}
	return
}
