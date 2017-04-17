package autorestart

import (
	"testing"

	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/slayer/autorestart"
	"github.com/stretchr/testify/assert"
)

func TestCustomFilenameAndFunc(t *testing.T) {
	filename := fmt.Sprintf("/tmp/autorestart-%d-%d", time.Now().Nanosecond(), os.Getpid())
	t.Logf("filename: %s", filename)
	err := ioutil.WriteFile(filename, []byte("1"), os.ModePerm)
	assert.Nil(t, err)

	customFuncOk := make(chan bool)
	autorestart.WatchFilename = filename
	autorestart.WatchPeriod = time.Millisecond
	autorestart.RestartFunc = func() {
		t.Logf("RestartFunc is fired")
		customFuncOk <- true
	}
	autorestart.StartWatcher()
	time.Sleep(time.Second)

	ioutil.WriteFile(filename, []byte("2"), os.ModePerm)

	select {
	case <-customFuncOk:
	case <-time.After(time.Second * 3):
		t.Errorf("Custom RestartFunc is not fired")
	}
	os.Remove(filename)
}
