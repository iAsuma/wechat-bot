package wechatbot

var botUserNickName = ""

func SetBotNickName(nickname string) {
	botUserNickName = nickname
}

func GetBotNickName() string {
	return botUserNickName
}
