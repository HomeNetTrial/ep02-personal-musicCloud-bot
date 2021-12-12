package musiccore

import (
	"errors"
	"fmt"
	"io/ioutil"
	"musicCloud-bot/log"
	"musicCloud-bot/pushmusic"
	"os"
	"strings"
)

type biliHandler struct{}

func (h biliHandler) IgetAudio(url string, infoChan chan<- string) (string, error) {
	if !Exists(gCacheFoldName) {
		os.Mkdir(gCacheFoldName, 0755)
	}
	tempName := makeTempID(url)
	sendMsg2chan("好呢,小🐷手收到Bilibili任务,正在处理中...", infoChan)
	checkPath := gCacheFoldName + "/" + tempName
	if Exists(checkPath) {
		files, err := ioutil.ReadDir(checkPath)
		if err == nil {
			for _, file := range files {
				if strings.Contains(file.Name(), ".mp3") {
					song := file.Name()
					fmt.Println("list files:" + song)
					sendMsg2chan("你可爱的小🐷手检查发现，此前已经处理过这个链接了,直接使用本地文件就行了哦:"+song, infoChan)
					return checkPath + "/" + song, nil
				}
			}
		}
	}
	mp3RealName, err := h.getVideoFile(url, tempName)
	if err != nil {
		return "", err
	}
	sendMsg2chan("小🐷手帮你获取到文件名称是:"+mp3RealName, infoChan)
	err = h.makeAudioFile(tempName, mp3RealName)
	if err == nil {
		tempPath := gCacheFoldName + "/" + tempName
		if !Exists(tempPath) {
			os.Mkdir(tempPath, 0755)
		}
		var param = []string{mp3RealName, tempPath}
		console, bRet := ExecCommand("mv", param)
		if !bRet {
			return "", errors.New(console)
		}
		mp3RealName = tempPath + "/" + mp3RealName
		var delParam = []string{tempName + ".mp4"}
		ExecCommand("rm", delParam)
	}
	return mp3RealName, err
}

//接口-把音频文件推送到网易云或者其他地方
func (h biliHandler) IpushAudio(path string, url string, infoChan chan<- string) error {
	err := pushmusic.PostFile(path, url)
	if err == nil {
		sendMsg2chan("你可爱的小🐷手已经帮你把这首歌["+path+"]推送成功了哦,快登陆一下【网易云音乐🎵】个人云盘瞅瞅!", infoChan)
	}
	return err
}
func (h biliHandler) getVideoFile(url string, tempFile string) (string, error) {
	cmdstr := "./annie"
	var param = []string{"-O", tempFile, url}
	consoleData, err := ExecCommand(cmdstr, param)
	if !err {
		return consoleData, errors.New(consoleData)
	}
	realName := h.getMp3FileName(consoleData)
	return realName, nil
}

func (h biliHandler) getMp3FileName(linesData string) string {
	start := strings.Index(linesData, "Title")
	end := strings.Index(linesData, "Type")
	//fileName := string([]rune(linesData)[start:end])
	fileName := linesData[start+len("Title:") : end]
	str := fileName + ".mp3"
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	log.Println(str)
	return str
}

func (h biliHandler) makeAudioFile(mp4Name string, mp3RealName string) error {
	cmdstr := "./ffmpeg"
	var param = []string{"-i", mp4Name + ".mp4", "-vn", "-f", "mp3", mp3RealName}
	console, err := ExecCommand(cmdstr, param)
	if !err {
		return errors.New(console)
	}
	return nil
}
