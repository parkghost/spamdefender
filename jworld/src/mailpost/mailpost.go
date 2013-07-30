package mailpost

import (
	"errors"
	"mailfile"
	"post"
	"post/article"
	"post/general"
	"strings"
)

var ErrNotSupport = errors.New("cannot parse non-JWorld@TW mails")

func Parse(mail mailfile.Mail) (*post.Post, error) {

	//JavaTWO即 將於 7/20盛 大舉辦！即日起可享早鳥優惠, 將New iPad帶 回家！
	if !strings.HasPrefix(mail.Subject(), "JWorld@TW") {
		return nil, ErrNotSupport
	}

	//JWorld@TW新文章通知:Hello JavaFX! Part 3
	for _, prefix := range []string{"JWorld@TW新文章通知:"} {
		if strings.HasPrefix(mail.Subject(), prefix) {
			return article.Parse(mail.Content())
		}
	}

	//JWorld@TW, 回覆通知
	if strings.HasPrefix(mail.Subject(), "JWorld@TW, 回覆通知") {
		return general.Parse(mail.Content())
	}

	//JWorld@TW新話題通知:*.class如何反編譯回*.java還能正常執行?
	//JWorld@TW話題更新通知:Guava 教學系列文章
	for _, prefix := range []string{"JWorld@TW新話題通知:", "JWorld@TW話題更新通知:", "JWorld@TW話題更新通知:"} {
		if strings.HasPrefix(mail.Subject(), prefix) {
			return general.Parse(mail.Content())
		}
	}

	return nil, ErrNotSupport
}
