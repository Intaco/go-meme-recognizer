package fetcher

import (
	"errors"
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

var package_logger = log.New(os.Stderr, "", 0)

var (
	ErrJobTimedOut        = errors.New("job request timed out")
	ErrConflictingDirName = errors.New("file named as provided dir name exist")
)

type Pull struct {
	concurrency int
}

type Task interface {
	Execute()
}

type Pool struct {
	mu    sync.Mutex
	size  int
	tasks chan Task
	kill  chan struct{}
	wg    sync.WaitGroup
}

const DEFAULT_POOL_TASK_CHANEL_SIZE = 128

func NewPool(size int) *Pool {
	pool := &Pool{
		tasks: make(chan Task, DEFAULT_POOL_TASK_CHANEL_SIZE),
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

func (p *Pool) Close() {
	close(p.tasks)
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Exec(e Task) {
	p.tasks <- e
}

type DownloadTask struct {
	query              DownloadQuery
	download_wait_time int
	logger             *log.Logger
}

func NewDownloadTask(query DownloadQuery, download_wait_time int,
	verbose int) (t *DownloadTask) {
	t = new(DownloadTask)
	t.query = query
	t.download_wait_time = download_wait_time
	if verbose > 0 {
		t.logger = log.New(os.Stderr, "DownloadTask: ", 0)
	}
	return
}

func (t *DownloadTask) download() (err error) {
	if t.logger != nil {
		t.logger.Printf("Started downloading %s\n\tto %s\n", t.query.url, t.query.getFilePath())
	}

	fd, err := t.query.Prepare()
	if err != nil {
		if t.logger != nil {
			t.logger.Println(err)
		}
		return
	}

	timeout := time.Duration(t.download_wait_time) * time.Second
	client := http.Client{Timeout: timeout}

	response, err := client.Get(t.query.url)
	if err != nil {
		if t.logger != nil {
			t.logger.Println(err)
		}
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(fd, response.Body)
	if err != nil {
		if t.logger != nil {
			t.logger.Println(err)
		}
		return
	}

	if t.logger != nil {
		t.logger.Printf("Finished downloading %s\n\t%s\n", t.query.url,
			t.query.getFilePath())
	}
	t.query.markAsDone()
	return
}

func (t DownloadTask) Execute() {
	t.download()
}

type DownloadQuery struct {
	url      string
	dirPath  string
	filename string
	isDone   bool
}

func NewDownloadQuery(url string, dirPath string,
	filename string) DownloadQuery {

	return DownloadQuery{url, dirPath, filename, false}
}

func (q *DownloadQuery) getFilePath() (filePath string) {
	filePath = filepath.Join(q.dirPath, q.filename)
	return
}

func (q *DownloadQuery) IsDone() bool {
	return q.isDone
}

func (q *DownloadQuery) markAsDone() {
	q.isDone = true
}

func (q *DownloadQuery) Prepare() (fd *os.File, err error) {
	if err = os.MkdirAll(q.dirPath, 0777); err != nil {
		package_logger.Println(err)
		return
	}
	filePath := q.getFilePath()
	if fd, err = os.Create(filePath); err != nil {
		package_logger.Println(err)
		return
	}

	return
}

type Fetcher struct {
	concurrency int
	wait_time   int
}

func NewFetcher(concurrency int, wait_time int) (f *Fetcher) {
	f = &Fetcher{concurrency, wait_time}
	return
}

func (f *Fetcher) download(downloads []DownloadQuery, startIndex int) {
	pool := NewPool(f.concurrency)
	for _, query := range downloads {
		pool.Exec(NewDownloadTask(query,
			f.wait_time, 1))
	}
	pool.Close()
	pool.Wait()
}

func makeQueryFromUrlsList(rootdir string, urlsList [][]string) (downloads []DownloadQuery) {
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

/*
func main() {
	fetcher := NewFetcher(10, 10)
	urlsList := [][]string{[]string{"https://pp.userapi.com/c841030/v841030005/1826e/Bunv2Om-uv4.jpg"},
		[]string{"https://pp.userapi.com/c638221/v638221662/81171/lLsKjoP3s_E.jpg",
			"https://pp.userapi.com/c638221/v638221602/55c39/TQnoSaS1eVI.jpg",
			"https://pp.userapi.com/c638221/v638221602/55c41/nOmwuC6O2mA.jpg",
			"https://pp.userapi.com/c638221/v638221602/55c48/g0J54ofxdyY.jpg",
			"https://pp.userapi.com/c638221/v638221602/55c4f/HAipw-io3uY.jpg"},
		[]string{"https://pp.userapi.com/c638221/v638221007/58c3c/9A0Tz4d06bc.jpg"},
		[]string{"https://pp.userapi.com/c638221/v638221007/58c32/nuAr6pMJGhs.jpg"},
		[]string{"https://pp.userapi.com/c638221/v638221007/58c04/KMyDz0wwDIc.jpg"},
		[]string{"https://pp.userapi.com/c638221/v638221388/5f8a1/_y7dUsi15b8.jpg"},
		[]string{"https://pp.userapi.com/c837731/v837731337/55cb3/yAsTav_Ap8A.jpg"},
		[]string{"https://pp.userapi.com/c837731/v837731869/5a3ce/0C9xZypRHRo.jpg"}}
	queries := makeQueryFromUrlsList("result", urlsList)
	fetcher.download(queries, 1)
}
*/
