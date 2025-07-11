package notification

import (
	"log"

	appConfig "github.com/akinoccc/hysaif/api/config"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/work/config"
)

// WeComMessageBuilder 企业微信消息构建器
type WeComMessageBuilder struct {
	title       string
	contentList []WeComContentItem
	jumpURL     string
	jumpTitle   string
}

// WeComContentItem 企业微信消息内容项
type WeComContentItem struct {
	KeyName string `json:"keyname"`
	Value   string `json:"value"`
}

// NewWeComMessageBuilder 创建新的企业微信消息构建器
func NewWeComMessageBuilder() *WeComMessageBuilder {
	return &WeComMessageBuilder{
		contentList: make([]WeComContentItem, 0),
	}
}

// SetTitle 设置消息标题
func (b *WeComMessageBuilder) SetTitle(title string) *WeComMessageBuilder {
	b.title = title
	return b
}

// AddContent 添加内容项
func (b *WeComMessageBuilder) AddContent(keyName, value string) *WeComMessageBuilder {
	b.contentList = append(b.contentList, WeComContentItem{
		KeyName: keyName,
		Value:   value,
	})
	return b
}

// SetJump 设置跳转链接
func (b *WeComMessageBuilder) SetJump(title, url string) *WeComMessageBuilder {
	b.jumpTitle = title
	b.jumpURL = url
	return b
}

// Build 构建企业微信消息体
func (b *WeComMessageBuilder) Build() map[string]interface{} {
	// 构建水平内容列表
	horizontalContentList := make([]interface{}, 0, len(b.contentList))
	for _, item := range b.contentList {
		horizontalContentList = append(horizontalContentList, map[string]interface{}{
			"keyname": item.KeyName,
			"value":   item.Value,
		})
	}

	// 构建模板卡片
	templateCard := map[string]interface{}{
		"card_type": "text_notice",
		"source": map[string]interface{}{
			"icon_url":   appConfig.AppConfig.Server.Domain + "/logo.png",
			"desc":       "Hysaif",
			"desc_color": 0,
		},
		"main_title": map[string]interface{}{
			"title": b.title,
		},
		"horizontal_content_list": horizontalContentList,
	}

	// 添加跳转链接（如果设置了）
	if b.jumpURL != "" {
		jumpList := []map[string]interface{}{
			{
				"type":  1,
				"url":   b.jumpURL,
				"title": b.jumpTitle,
			},
		}
		templateCard["jump_list"] = jumpList
		templateCard["card_action"] = map[string]interface{}{
			"type": 1,
			"url":  b.jumpURL,
		}
	}

	return map[string]interface{}{
		"msgtype":       "template_card",
		"template_card": templateCard,
	}
}

// Send 发送消息
func (b *WeComMessageBuilder) Send() error {
	msgBody := b.Build()
	return SendWeComMessage(msgBody)
}

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

	info, err := wc.GetRobot().RobotBroadcast(
		appConfig.AppConfig.WeCom.RobotHookKey,
		msgBody,
	)
	log.Println("wecom:===============", info.ErrMsg)
	return err
}

// 便捷方法：发送访问申请通知
func SendAccessRequestNotification(title, applicantName, secretItemName, reason, requestID string) error {
	return NewWeComMessageBuilder().
		SetTitle(title).
		AddContent("申请人", applicantName).
		AddContent("密钥项", secretItemName).
		AddContent("申请理由", reason).
		SetJump("查看申请", appConfig.AppConfig.Server.Domain+"/access_requests?id="+requestID).
		Send()
}

// 便捷方法：发送访问申请批准通知
func SendAccessApprovedNotification(approverName, secretItemName, validUntil, note, secretItemID, secretItemType string) error {
	return NewWeComMessageBuilder().
		SetTitle("敏感信息审批申请已批准").
		AddContent("审批人", approverName).
		AddContent("密钥项", secretItemName).
		AddContent("有效期至", validUntil).
		AddContent("备注", note).
		SetJump("查看密钥项", appConfig.AppConfig.Server.Domain+"/"+secretItemType+"?id="+secretItemID).
		Send()
}

// 便捷方法：发送访问申请拒绝通知
func SendAccessRejectedNotification(approverName, secretItemName, rejectReason, requestID string) error {
	return NewWeComMessageBuilder().
		SetTitle("敏感信息审批申请已拒绝").
		AddContent("审批人", approverName).
		AddContent("密钥项", secretItemName).
		AddContent("拒绝原因", rejectReason).
		SetJump("查看申请", appConfig.AppConfig.Server.Domain+"/access_requests?id="+requestID).
		Send()
}

// 便捷方法：发送访问申请作废通知
func SendAccessRevokedNotification(operatorName, secretItemName, revokeReason, requestID string) error {
	return NewWeComMessageBuilder().
		SetTitle("敏感信息审批申请已作废").
		AddContent("操作人", operatorName).
		AddContent("密钥项", secretItemName).
		AddContent("作废原因", revokeReason).
		SetJump("查看申请", appConfig.AppConfig.Server.Domain+"/access_requests?id="+requestID).
		Send()
}

// 便捷方法：发送访问申请过期通知
func SendAccessExpiredNotification(secretItemName, requestID string) error {
	return NewWeComMessageBuilder().
		SetTitle("敏感信息审批申请已过期").
		AddContent("密钥项", secretItemName).
		SetJump("查看申请", appConfig.AppConfig.Server.Domain+"/access_requests?id="+requestID).
		Send()
}

// 便捷方法：发送密钥项即将过期通知
func SendSecretItemExpiringNotification(secretItemName, expiresIn, secretItemID, secretItemType string) error {
	return NewWeComMessageBuilder().
		SetTitle("密钥项即将过期").
		AddContent("密钥项", secretItemName).
		AddContent("过期时间", expiresIn).
		SetJump("查看密钥项", appConfig.AppConfig.Server.Domain+"/"+secretItemType+"?id="+secretItemID).
		Send()
}

// 便捷方法：发送自定义通知
func SendCustomNotification(title string, contents map[string]string, jumpTitle, jumpURL string) error {
	builder := NewWeComMessageBuilder().SetTitle(title)

	// 添加所有内容项
	for key, value := range contents {
		builder.AddContent(key, value)
	}

	// 设置跳转链接（如果提供）
	if jumpURL != "" {
		builder.SetJump(jumpTitle, jumpURL)
	}

	return builder.Send()
}

/*
使用示例：

1. 基本用法（使用便捷方法）:
err := SendAccessRequestNotification(
    "新的访问申请",
    "张三",
    "生产环境数据库密码",
    "需要紧急修复线上问题",
    "req-123",
)

2. 高级用法（使用 Builder）:
err := NewWeComMessageBuilder().
    SetTitle("系统维护通知").
    AddContent("维护时间", "2024-01-01 02:00-04:00").
    AddContent("影响范围", "用户登录功能").
    AddContent("预计时长", "2小时").
    SetJump("查看详情", "https://status.example.com").
    Send()

3. 自定义通知:
err := SendCustomNotification(
    "安全警报",
    map[string]string{
        "检测时间": "2024-01-01 10:30:00",
        "异常IP": "192.168.1.100",
        "尝试次数": "5次",
    },
    "查看日志",
    "https://log.example.com/security",
)

4. 仅构建消息体（不发送）:
msgBody := NewWeComMessageBuilder().
    SetTitle("测试消息").
    AddContent("环境", "开发环境").
    Build()
// 然后可以用于其他用途，如保存到数据库或传递给其他服务
*/
