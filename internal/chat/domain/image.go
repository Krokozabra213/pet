package chatdomain

type IImage interface {
	SetMessageID(id uint64)
	GetImageUrl() string
	GetMessageID() uint64
}

type Image struct {
	ID        uint64  `gorm:"primarykey;column:id"`
	MessageID uint64  `gorm:"uniqueIndex;column:message_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ImageUrl  string  `gorm:"column:image_url"`
	Message   Message `gorm:"foreignKey:MessageID;references:ID"`
}

func NewImage(url string) *Image {
	return &Image{
		ImageUrl: url,
	}
}

func (i *Image) SetMessageID(id uint64) {
	i.MessageID = id
}

func (i *Image) GetMessageID() uint64 {
	return i.MessageID
}

func (i *Image) GetImageUrl() string {
	return i.ImageUrl
}
