package emailservice

import (
	"fmt"

	"github.com/Monkey-Mouse/mo2/mo2utils/mo2errors"
)

// VerifyEmailMessage build email body for verify email
func VerifyEmailMessage(url string, name string, receivers []string) *Mo2Email {

	// subject := "Subject: 确认Mo2邮箱!\r\n"
	// mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	body := "<html><body>" +
		"  <tr>\n        <td>\n            <div style=\"border: #36649d 1px dashed;margin: 30px;padding: 20px\">\n                <div style=\"margin: 40px\">\n                    <label style=\"font-size: 22px;color: #36649d;font-weight: bold\">请确认你的邮箱！</label>\n                    <p style=\"font-size: 16px\">亲爱的&nbsp;<label style=\"font-weight: bold\"> " +
		name + "</label>&nbsp; 您好！欢迎来到Mo2\n                    </p>\n                    <p style=\"font-size: 16px\">在开始使用Mo2之前，您必须先确认您的电子邮件地址</p>\n                </div>\n\n                <div class=\"col-pad-left-3 col-pad-right-2\" style=\"color:#0a0a0a;font-family:'Cereal', Helvetica, Arial, sans-serif;font-weight:normal;padding:40px;margin:0;text-align:center;font-size:16px;line-height:19px;padding-left:16px;padding-right:16px\">\n        " +
		"            <a href=\"" +
		url + "\"  class=\"btn-primary btn-md btn-rausch\" style=\"font-family:'Cereal', Helvetica, Arial, sans-serif;font-weight:normal;margin:0;text-align:left;line-height:1.3;color:#2199e8;text-decoration:none;background-color:#ff5a5f;-webkit-border-radius:4px;border-radius:4px;display:inline-block;padding:12px 24px 12px 24px;\" rel=\"noopener\" target=\"_blank\">\n                        <p class=\"text-center\" style=\"font-weight:normal;padding:0;margin:0;text-align:center;font-family:&quot;Cereal;      &quot;, &quot;Helvetica&quot;, Helvetica, Arial, sans-serif;color:white;font-size:18px;line-height:26px;margin-bottom:0px !important;\">\n                            确认邮箱\n                        </p>\n                    </a>\n                </div>\n                <div style=\"margin: 40px\">\n                    <p style=\"font-size: 16px\">谢谢！</p>\n                    <p style=\"font-size: 16px\">Mo2团队</p>\n                    <p style=\"color:red;font-size: 14px \">（这是一封自动发送的邮件，请勿回复）</p>\n\n                </div>\n            </div>\n        </td>\n\n    </tr></body></html>"

	return &Mo2Email{Content: body, Subject: "确认Mo2邮箱!", Receivers: receivers}
}

// InvitationMessage build email body for invite friends
func InvitationMessage(url, name, groupname string, receivers []string) *Mo2Email {

	// subject := "Subject: 确认Mo2邮箱!\r\n"
	// mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	body := "<html><body>" +
		"  <tr>\n        <td>\n            <div style=\"border: #36649d 1px dashed;margin: 30px;padding: 20px\">\n                <div style=\"margin: 40px\">\n                    <label style=\"font-size: 22px;color: #36649d;font-weight: bold\">Group邀请</label>\n                    <p style=\"font-size: 16px\">亲爱的&nbsp;<label style=\"font-weight: bold\"> " +
		name + "</label>&nbsp; 您好！您收到一封来自Group" + groupname + "的邀请！  \n                  </p>\n                    <p style=\"font-size: 16px\">点击下方按钮即可加入</p>\n                </div>\n\n                <div class=\"col-pad-left-3 col-pad-right-2\" style=\"color:#0a0a0a;font-family:'Cereal', Helvetica, Arial, sans-serif;font-weight:normal;padding:40px;margin:0;text-align:center;font-size:16px;line-height:19px;padding-left:16px;padding-right:16px\">\n        " +
		"            <a href=\"" +
		url + "\"  class=\"btn-primary btn-md btn-rausch\" style=\"font-family:'Cereal', Helvetica, Arial, sans-serif;font-weight:normal;margin:0;text-align:left;line-height:1.3;color:#2199e8;text-decoration:none;background-color:#ff5a5f;-webkit-border-radius:4px;border-radius:4px;display:inline-block;padding:12px 24px 12px 24px;\" rel=\"noopener\" target=\"_blank\">\n                        <p class=\"text-center\" style=\"font-weight:normal;padding:0;margin:0;text-align:center;font-family:&quot;Cereal;      &quot;, &quot;Helvetica&quot;, Helvetica, Arial, sans-serif;color:white;font-size:18px;line-height:26px;margin-bottom:0px !important;\">\n                            确认接受\n                        </p>\n                    </a>\n                </div>\n                <div style=\"margin: 40px\">\n                    <p style=\"font-size: 16px\">谢谢！</p>\n                    <p style=\"font-size: 16px\">Mo2团队</p>\n                    <p style=\"color:red;font-size: 14px \">（这是一封自动发送的邮件，请勿回复）</p>\n\n                </div>\n            </div>\n        </td>\n\n    </tr></body></html>"

	return &Mo2Email{Content: body, Subject: "Group邀请", Receivers: receivers}
}

// SendEmail can send raw email to receivers
func SendEmail(mail *Mo2Email, senderAddr string) (err *mo2errors.Mo2Errors) {

	err = QueueEmail(mail, senderAddr)
	if err != nil {
		fmt.Println(err)
	}
	return
}
