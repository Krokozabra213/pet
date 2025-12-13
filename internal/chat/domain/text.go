package chatdomain

type IText interface {
	GetContent() string
}

type Text struct {
	id        uint64 `gorm:"primarykey;column:id"`
	messageID uint64 `gorm:"uniqueIndex;column:message_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	content   string `gorm:"column:content"`
}

func NewText(content string) *Text {
	return &Text{
		content: content,
	}
}

func (t *Text) GetContent() string {
	return t.content
}

func (t *Text) GetMessageID() uint64 {
	return t.messageID
}

func (t *Text) SetMessageID(id uint64) {
	t.messageID = id
}
