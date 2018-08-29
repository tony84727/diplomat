package diplomat

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type IBuildSpace interface {
	Create(path string) (io.WriteCloser, error)
}

type BuildSpace struct {
	dir    string
	logger *log.Logger
	prefix string
	wg     sync.WaitGroup
}

func newBuildSpace(dir string) BuildSpace {
	return BuildSpace{
		dir:    dir,
		logger: log.New(os.Stdout, "", log.Flags()),
	}
}

func (b BuildSpace) ForMessenger(messengerType string) MessengerBuildSpace {
	mbs := MessengerBuildSpace{b}
	mbs.setPrefix(fmt.Sprintf("[%s]", messengerType))
	return mbs
}

type spyWriteCloser struct {
	target    io.WriteCloser
	couter    chan<- int
	sum       int
	listeners []func(int)
}

func newspyWriteCloser(target io.WriteCloser) *spyWriteCloser {
	couterChan := make(chan int, 1)
	sw := &spyWriteCloser{
		target:    target,
		couter:    couterChan,
		sum:       0,
		listeners: make([]func(int), 0, 1),
	}
	go func() {
		for i := range couterChan {
			sw.sum += i
		}
	}()
	return sw
}

func (sw *spyWriteCloser) Write(p []byte) (n int, err error) {
	n, err = sw.target.Write(p)
	sw.couter <- n
	return
}

func (sw *spyWriteCloser) Close() error {
	err := sw.target.Close()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, listener := range sw.listeners {
		wg.Add(1)
		go func(f func(int)) {
			f(sw.sum)
			wg.Done()
		}(listener)
	}
	wg.Wait()
	return err
}

func (sw *spyWriteCloser) addCloseListener(listener func(int)) {
	sw.listeners = append(sw.listeners, listener)
}

func (b BuildSpace) Println(content ...interface{}) {
	content = append([]interface{}{b.prefix}, content...)
	b.logger.Println(content...)
}

func (b BuildSpace) onWriteFileDone(path string, size int) {
	b.Println(fmt.Sprintf("[emit] %d bytes to %s", size, path))
}

func (b BuildSpace) Create(path string) (io.WriteCloser, error) {
	path = filepath.Join(b.dir, path)
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	b.wg.Add(1)
	spy := newspyWriteCloser(f)
	spy.addCloseListener(func(size int) {
		b.onWriteFileDone(path, size)
		b.wg.Done()
	})
	return spy, nil
}

func (b BuildSpace) Done() {
	b.wg.Wait()
}

func (b *BuildSpace) setPrefix(prefix string) {
	b.prefix = prefix
}

type MessengerBuildSpace struct {
	BuildSpace
}
