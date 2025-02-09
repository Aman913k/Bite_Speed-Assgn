package models

import "time"

// Contact model
type Contact struct {
	ID             int        `json:"id"`
	PhoneNumber    *string    `json:"phoneNumber,omitempty"`
	Email          *string    `json:"email,omitempty"`
	LinkedID       *int       `json:"linkedId,omitempty"`
	LinkPrecedence string     `json:"linkPrecedence"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty"`
}

// Response model
type IdentifyResponse struct {
	Contact struct {
		PrimaryContactID    int      `json:"primaryContactId"`
		Emails              []string `json:"emails"`
		PhoneNumbers        []string `json:"phoneNumbers"`
		SecondaryContactIDs []int    `json:"secondaryContactIds"`
	} `json:"contact"`
}
