package chatdomain

type IImage interface {
	SetMessageID(id uint64)
	GetImageUrl() string
	GetMessageID() uint64
}

type Image struct {
	id        uint64 `gorm:"primarykey;column:id"`
	messageID uint64 `gorm:"uniqueIndex;column:message_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	imageUrl  string `gorm:"column:image_url"`
}

func NewImage(url string) *Image {
	return &Image{
		imageUrl: url,
	}
}

func (i *Image) SetMessageID(id uint64) {
	i.messageID = id
}

func (i *Image) GetMessageID() uint64 {
	return i.messageID
}

func (i *Image) GetImageUrl() string {
	return i.imageUrl
}
