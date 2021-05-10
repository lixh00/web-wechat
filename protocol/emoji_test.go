package protocol

import (
	"testing"
)

func TestFormatEmoji(t *testing.T) {
	t.Log(FormatEmoji(`这是一个苹果<span class="emoji emoji1f34f"></span>`))
	t.Log(FormatEmoji(`这里没有苹果`))
	t.Log(FormatEmoji(""))
}

func TestSendEmoji(t *testing.T) {
	self, err := getSelf()
	if err != nil {
		t.Error(err)
		return
	}
	f, err := self.FileHelper()
	if err != nil {
		t.Error(err)
		return
	}
	_, _ = f.SendText(Emoji.Dagger)
}
