package musiccore

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"musicCloud-bot/config"
	"net/http"
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func getRealUrl(url string) (string, error) {

	if strings.Contains(url, "bilibili.com") || strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
		return url, nil
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get(url)
	if err != nil {
		return url, err
	}
	fmt.Println(res)
	fmt.Println("StatusCode:" + strconv.Itoa(res.StatusCode))
	if (res.StatusCode != 301) && (res.StatusCode != 302) {
		return url, nil
	}
	Location := res.Header.Get("Location")
	fmt.Println("res.Header.Get(\"Location\"):", Location)
	return Location, nil
}

type reqPack struct {
	taskID    string
	t         int64
	recvMsg   string
	message   *tb.Message
	sendMsgCB func(m *tb.Message, msg string)
	h         Ihandler
}

func NewReqPack(Msg string, m *tb.Message, f func(m *tb.Message, msg string)) reqPack {
	obj := reqPack{
		t:         time.Now().Unix(),
		taskID:    "",
		recvMsg:   Msg,
		sendMsgCB: f,
		message:   m,
		h:         nil,
	}
	return obj
}

func (r reqPack) HandleMsg() {
	handle := r.checkAndGetHandler(r.recvMsg)
	if handle == nil {
		r.sendMsgCB(r.message, "小🐷手暂时只支持哔哩哔哩的视频链接哦，YouTube的视频正在计划中呢!等我哦！")
		return
	}
	r.h = handle
	r.taskID = r.mkaeTaskID(r.recvMsg)
	ctx, cancle := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancle()
	infoChan := make(chan string)
	errChan := make(chan error)

	//do get the video to mp3 task
	go func(ctx context.Context, r reqPack, infoC chan<- string) {
		path, err := r.h.IgetAudio(r.recvMsg, infoChan)
		if err != nil {
			r.sendBackMsg(r.taskID, "获取音频文件失败了")
			errChan <- err
		}
		errChan <- r.h.IpushAudio(path, config.PushMusicServerHost, infoChan)
		<-ctx.Done()
		useTime := time.Now().Unix() - r.t
		r.sendBackMsg(r.taskID, "任务已经被结束了,这次用时:"+strconv.FormatInt(useTime, 10)+"秒!")
	}(ctx, r, infoChan)
	//send running log
	go func(ctx context.Context, r reqPack, infoC chan string) {
		for {
			select {
			case msg := <-infoC:
				r.sendBackMsg(r.taskID, msg)
			case <-ctx.Done():
				return
			}
		}
	}(ctx, r, infoChan)

	err := <-errChan
	if err != nil {
		print := "任务捕捉到错误信息:" + err.Error()
		log.Println(print)
		r.sendBackMsg(r.taskID, print)
	}
}

func (r reqPack) sendBackMsg(taskId string, msg string) {
	preFix := "[" + time.Now().Format("15:04:05") + "任务编号:" + taskId + "]"
	r.sendMsgCB(r.message, preFix+msg)
}

func (r reqPack) mkaeTaskID(url string) string {
	index := strings.Index(url, "?")
	if -1 != index {
		url = url[0:index]
	}
	data := []byte(url)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str[0:4]
}

func (r *reqPack) checkAndGetHandler(url string) Ihandler {
	r.recvMsg, _ = getRealUrl(url)
	url = r.recvMsg
	if strings.Contains(url, "bilibili.com") {
		index := strings.Index(url, "?")
		if -1 != index {
			url = url[0:index]
		}
		r.recvMsg = url
		return biliHandler{}
	} else if strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
		index := strings.Index(url, "&")
		if -1 != index {
			url = url[0:index]
		}
		r.recvMsg = url
		return youtubeHandler{}
	}
	return nil
}
