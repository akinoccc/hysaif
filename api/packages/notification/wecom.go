package notification

import (
	appConfig "github.com/akinoccc/hysaif/api/config"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/work/config"
)

// SendWeComMessage 发送企业微信消息
func SendWeComMessage(msgBody interface{}) error {
	if !appConfig.AppConfig.WeCom.Enabled {
		return nil
	}

	cacher := cache.NewMemory()

	wc := wechat.NewWechat().GetWork(&config.Config{
		CorpID:     appConfig.AppConfig.WeCom.CorpID,
		AgentID:    appConfig.AppConfig.WeCom.AgentID,
		CorpSecret: appConfig.AppConfig.WeCom.Secret,
		Cache:      cacher,
	})

	// _, err := wc.GetMessage().SendText(message.SendTextRequest{
	// 	SendRequestCommon: &message.SendRequestCommon{
	// 		ToUser:                 "@all",
	// 		AgentID:                appConfig.AppConfig.WeCom.AgentID,
	// 		Safe:                   1, // 0: 可对外分享，1: 不可分享且内容显示水印
	// 		EnableDuplicateCheck:   0,
	// 		DuplicateCheckInterval: 1800,
	// 	},
	// 	Text: message.TextField{
	// 		Content: msg,
	// 	},
	// })

	_, err := wc.GetRobot().RobotBroadcast(
		appConfig.AppConfig.WeCom.RobotHookKey,
		msgBody,
	)
	return err
}
