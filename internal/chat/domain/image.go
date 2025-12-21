package chatdomain

// type IImage interface {
// 	SetMessageID(id uint64)
// 	GetImageUrl() string
// 	GetMessageID() uint64
// }

type Image struct {
	ID        uint64 `gorm:"column:id"`
	MessageID uint64 `gorm:"column:message_id"`
	ImageUrl  string `gorm:"column:image_url"`
}

func NewImage(url string) *Image {
	return &Image{
		ImageUrl: url,
	}
}

func (i *Image) TableName() string {
	return "images"
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
