package pushmusic

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	gResp_login *http.Response
	loginUrl    string
)

func Login(serverHost string, user string, passwd string) (bool, error) {
	gResp_login = nil
	loginUrl = serverHost + "/login/cellphone?phone=" + user + "&md5_password=" + passwd
	resp_login, err := http.Get(loginUrl)
	if err != nil {
		return false, err
	}
	gResp_login = resp_login
	return true, nil
}
func reLogin() (bool, error) {
	resp_login, err := http.Get(loginUrl)
	if err != nil {
		return false, err
	}
	gResp_login = resp_login
	return true, nil
}

func PostFile(filename string, serverHost string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	now := time.Now().Unix()
	targetUrl := serverHost + "/cloud?time=" + strconv.FormatInt(now, 10)
	// 关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("songFile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// 打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	// iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	//先登录系统，再获取cookie
	// if gResp_login == nil {
	// 	ret, err := reLogin()
	// 	if !ret {
	// 		log.Println(err)
	// 		return err
	// 	}
	// }
	reLogin()

	req, err := http.NewRequest("POST", targetUrl, bodyBuf)
	if err != nil {
		fmt.Println("req err: ", err)
		return err
	}
	req.Header.Set("Content-Type", contentType)
	for _, c := range gResp_login.Cookies() {
		req.AddCookie(c)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("resp err: ", err)
		return err
	}

	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}
