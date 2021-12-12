package musiccore

import (
	"errors"
	"fmt"
	"io/ioutil"
	"musicCloud-bot/config"
	"musicCloud-bot/pushmusic"
	"os"
	"strings"
)

type youtubeHandler struct {
}

func (h youtubeHandler) IgetAudio(url string, infoChan chan<- string) (string, error) {
	if !Exists(gCacheFoldName) {
		os.Mkdir(gCacheFoldName, 0755)
	}
	tempName := makeTempID(url)
	sendMsg2chan("好呢,小🐷手收到Youtube任务,正在处理中...", infoChan)
	checkPath := gCacheFoldName + "/" + tempName
	if Exists(checkPath) {
		files, err := ioutil.ReadDir(checkPath)
		if err == nil {
			for _, file := range files {
				if strings.Contains(file.Name(), ".m4a") {
					song := file.Name()
					fmt.Println("list files:" + song)
					sendMsg2chan("你可爱的小🐷手检查发现，此前已经处理过这个链接了,直接使用本地文件就行了哦:"+song, infoChan)
					return checkPath + "/" + song, nil
				}
			}
		}
	} else {
		os.Mkdir(checkPath, 0755)
	}
	mp3RealName, err := h.makeAudioFile(url, checkPath)
	if err == nil {
		sendMsg2chan("小🐷手帮你获取到文件名称是:"+mp3RealName[len(checkPath)+1:], infoChan)
	}
	return mp3RealName, err
}

//接口-把音频文件推送到网易云或者其他地方
func (h youtubeHandler) IpushAudio(path string, url string, infoChan chan<- string) error {
	err := pushmusic.PostFile(path, url)
	if err == nil {
		sendMsg2chan("你可爱的小🐷手已经帮你把这首歌["+path+"]推送成功了哦,快登陆一下【网易云音乐🎵】个人云盘瞅瞅!", infoChan)
	}
	return err
}

func (h youtubeHandler) getAuduoFileName(cachePath string) string {
	files, err := ioutil.ReadDir(cachePath)
	if err == nil {
		for _, file := range files {
			if strings.Contains(file.Name(), ".m4a") {
				song := file.Name()
				return cachePath + "/" + song
			}
		}
	}
	return ""
}
func (h youtubeHandler) makeAudioFile(url string, path string) (string, error) {
	cmdstr := "youtube-dl"
	var param []string
	fileName := "%(title)s-%(id)s.m4a"
	if config.Socks5 != "" {
		proxyUrl := fmt.Sprintf("socks5://%s", config.Socks5)
		param = []string{"--proxy", proxyUrl, "-o", path + "/" + fileName, "-f", "140", url}
	} else {
		param = []string{"-o", path + "/" + fileName, "-f", "140", url}
	}

	console, err := ExecCommand(cmdstr, param)
	if !err {
		return "", errors.New(console)
	}
	file := h.getAuduoFileName(path)
	if file == "" {
		err := errors.New("getAuduoFileName fail:" + console)
		return "", err
	}
	return file, nil
}
