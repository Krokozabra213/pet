package newdomain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type School struct {
	ID          primitive.ObjectID
	Name        string
	Description string
	CreatedAt   time.Time
	Admins      []Admin
	Courses     []CourseSubCollection
	Info        Info
	Published   bool
}

type CourseSubCollection struct {
	ID          primitive.ObjectID
	Name        string
	Description string
	Color       string
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Info struct {
	Color               string
	Domains             []string
	Contacts            SchoolContacts
	Logo                string
	GoogleAnalyticsCode string
}

type SchoolContacts struct {
	BusinessName       string
	RegistrationNumber string
	Address            string
	Email              string
	Phone              string
}
