package domain

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
	Courses     []Course
	Info        Info
	Published   bool
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

type Admin struct {
	UserID   uint64
	Name     string
	Email    string
	SchoolID primitive.ObjectID
}
