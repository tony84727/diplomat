package diplomat

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type IBuildSpace interface {
	Create(path string) (io.WriteCloser, error)
}

type BuildSpace struct {
	dir    string
	logger *log.Logger
	prefix string
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
	for _, listener := range sw.listeners {
		go listener(sw.sum)
	}
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
	f, err := os.Create(filepath.Join(b.dir, path))
	if err != nil {
		return nil, err
	}
	spy := newspyWriteCloser(f)
	spy.addCloseListener(func(size int) {
		b.onWriteFileDone(path, size)
	})
	return spy, nil
}

func (b *BuildSpace) setPrefix(prefix string) {
	b.prefix = prefix
}

type MessengerBuildSpace struct {
	BuildSpace
}
