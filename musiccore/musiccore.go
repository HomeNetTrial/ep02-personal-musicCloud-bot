package musiccore

import (
	"crypto/md5"
	"fmt"
)

var (
	gCacheFoldName string = "cache"
)

type Ihandler interface {
	/* TODO: add methods */
	IgetAudio(url string, infoChan chan<- string) (string, error)
	IpushAudio(path string, url string, infoChan chan<- string) error
}

func sendMsg2chan(msg string, infoChan chan<- string) {
	infoChan <- msg
}

func makeTempID(url string) string {
	data := []byte(url)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str[0:12]
}
