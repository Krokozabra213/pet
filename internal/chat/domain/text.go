package chatdomain

// type IText interface {
// 	GetContent() string
// }

type Text struct {
	ID        uint64 `gorm:"column:id"`
	MessageID uint64 `gorm:"column:message_id"`
	Content   string `gorm:"column:content"`
}

func NewText(content string) *Text {
	return &Text{
		Content: content,
	}
}

func (t *Text) TableName() string {
	return "texts"
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
