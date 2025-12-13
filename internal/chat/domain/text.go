package chatdomain

type IText interface {
	GetContent() string
}

type Text struct {
	ID        uint64  `gorm:"primarykey;column:id"`
	MessageID uint64  `gorm:"uniqueIndex;column:message_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Content   string  `gorm:"column:content"`
	Message   Message `gorm:"foreignKey:MessageID;references:ID"`
}

func NewText(content string) *Text {
	return &Text{
		Content: content,
	}
}

func (t *Text) GetContent() string {
	return t.Content
}

func (t *Text) GetMessageID() uint64 {
	return t.MessageID
}

func (t *Text) SetMessageID(id uint64) {
	t.MessageID = id
}
