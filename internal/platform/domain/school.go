package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type School struct {
// 	ID          primitive.ObjectID
// 	Name        string
// 	Description string
// 	CreatedAt   time.Time
// 	Admins      []Admin
// 	Courses     []CourseSubCollection
// 	Info        Info
// 	Published   bool
// }

// type CourseSubCollection struct {
// 	ID          primitive.ObjectID
// 	Name        string
// 	Description string
// 	Color       string
// 	ImageURL    string
// 	CreatedAt   time.Time
// 	UpdatedAt   time.Time
// }

// type Info struct {
// 	Color               string
// 	Domains             []string
// 	Contacts            SchoolContacts
// 	Logo                string
// 	GoogleAnalyticsCode string
// }

// type SchoolContacts struct {
// 	BusinessName       string
// 	RegistrationNumber string
// 	Address            string
// 	Email              string
// 	Phone              string
// }

type School struct {
	ID          primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Name        string                `json:"name" bson:"name"`
	Description string                `json:"description" bson:"description"`
	CreatedAt   time.Time             `json:"created_at" bson:"created_at"`
	Admins      []Admin               `json:"admins" bson:"admins"`
	Courses     []CourseSubCollection `json:"courses" bson:"courses"`
	Info        Info                  `json:"info" bson:"info"`
	Published   bool                  `json:"published" bson:"published"`
}

type CourseSubCollection struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Color       string             `json:"color" bson:"color"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type Info struct {
	Color               string         `json:"color" bson:"color"`
	Domains             []string       `json:"domains" bson:"domains"`
	Contacts            SchoolContacts `json:"contacts" bson:"contacts"`
	Logo                string         `json:"logo" bson:"logo"`
	GoogleAnalyticsCode string         `json:"google_analytics_code" bson:"google_analytics_code"`
}

type SchoolContacts struct {
	BusinessName       string `json:"business_name" bson:"business_name"`
	RegistrationNumber string `json:"registration_number" bson:"registration_number"`
	Address            string `json:"address" bson:"address"`
	Email              string `json:"email" bson:"email"`
	Phone              string `json:"phone" bson:"phone"`
}

type SchoolCreateInput struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name                string             `json:"name" bson:"name"`
	Description         string             `json:"description" bson:"description"`
	Color               string             `json:"color" bson:"color"`
	Logo                string             `json:"logo" bson:"logo"`
	GoogleAnalyticsCode string             `json:"google_analytics_code" bson:"google_analytics_code"`
	BusinessName        string             `json:"business_name" bson:"business_name"`
	RegistrationNumber  string             `json:"registration_number" bson:"registration_number"`
	Address             string             `json:"address" bson:"address"`
	Email               string             `json:"email" bson:"email"`
	Phone               string             `json:"phone" bson:"phone"`
	Published           bool               `json:"published" bson:"published"`
}

type SchoolUpdateInput struct {
	Name                *string `json:"name" bson:"name"`
	Description         *string `json:"description" bson:"description"`
	Color               *string `json:"color" bson:"color"`
	Logo                *string `json:"logo" bson:"logo"`
	GoogleAnalyticsCode *string `json:"google_analytics_code" bson:"google_analytics_code"`
	BusinessName        *string `json:"business_name" bson:"business_name"`
	RegistrationNumber  *string `json:"registration_number" bson:"registration_number"`
	Address             *string `json:"address" bson:"address"`
	Email               *string `json:"email" bson:"email"`
	Phone               *string `json:"phone" bson:"phone"`
	Published           *bool   `json:"published" bson:"published"`
}

func (s SchoolUpdateInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	if s.Name != nil {
		m["name"] = *s.Name
	}
	if s.Description != nil {
		m["description"] = *s.Description
	}
	if s.Color != nil {
		m["color"] = *s.Color
	}
	if s.Logo != nil {
		m["logo"] = *s.Logo
	}
	if s.GoogleAnalyticsCode != nil {
		m["google_analytics_code"] = *s.GoogleAnalyticsCode
	}
	if s.BusinessName != nil {
		m["business_name"] = *s.BusinessName
	}
	if s.RegistrationNumber != nil {
		m["registration_number"] = *s.RegistrationNumber
	}
	if s.Address != nil {
		m["address"] = *s.Address
	}
	if s.Email != nil {
		m["email"] = *s.Email
	}
	if s.Phone != nil {
		m["phone"] = *s.Phone
	}
	if s.Phone != nil {
		m["phone"] = *s.Phone
	}
	if s.Published != nil {
		m["published"] = *s.Published
	}

	return m
}
