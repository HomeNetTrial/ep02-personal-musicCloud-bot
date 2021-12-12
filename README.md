## musicCloud-bot
### 自动化推送机器人

如果你像我一样，喜欢逛B站、油管，喜欢把一些UP主的作品转成MP3音频文件反复收听，可以尝试一下我写的这个机器人助手。目前只支持电报(telegram)的消息的发送。

![](https://cdn.jsdelivr.net/gh/playbear668/mypics//img/5691639290502_.pic.jpg)



### 源码的编译和依赖

- 目前只支持Telegram(电报)的机器人，如需要其他可以自行更改
- 推送到网易云音乐部分，使用以下开源项目,需先行部署该服务,感谢Binaryify大大开源这么棒的API
```shell
https://github.com/Binaryify/NeteaseCloudMusicApi
```
- 电报机器人部分代码采用以下开源项目修改
```shell
https://github.com/indes/flowerss-bot
```
- 需要使用到annie、ffmpeg、youtube-dl等开源下面的组件进行音视频的转换，非常感谢上述开源大佬的基础组件的支持！！

### 懒人福利

如果你实在不向自己编译源码出包，又想体验一下的话，可以赞赏一下我这边文章，请我吃颗糖，我送你个打包好的执行程序，开箱即用！！

不过暂时也只出了linux端的X86的软件包。像树莓派这样的arm架构，等看多人需要，再考虑出吧。


### 你的专属定制

如果开源代码也没办法满足你的需要，你有一些更不错的想法，也是支持付费定制的，请添加我的私信进行了解哦。

