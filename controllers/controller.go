// controllers.go
package controllers

import (
	"Bite-Speed/database"
	"Bite-Speed/models"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func IdentifyHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       *string `json:"email"`
		PhoneNumber *string `json:"phoneNumber"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == nil && req.PhoneNumber == nil {
		http.Error(w, "At least one of email or phoneNumber is required", http.StatusBadRequest)
		return
	}

	contacts, _ := FetchLinkedContacts(req.Email, req.PhoneNumber)

	var primaryID *int
	precedence := "primary"

	if len(contacts) > 0 {
		primary, _ := SegregateContacts(contacts)
		precedence = "secondary"
		primaryID = &primary.ID
	}

	// Inserting new contact
	newID := CreateContact(req.Email, req.PhoneNumber, primaryID, precedence)
	contacts = append(contacts, models.Contact{
		ID:             newID,
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		LinkPrecedence: precedence,
		LinkedID:       primaryID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	contacts, _ = FetchLinkedContacts(req.Email, req.PhoneNumber)
	response := FormatResponse(contacts)
	json.NewEncoder(w).Encode(response)
}

func FetchLinkedContacts(email, phoneNumber *string) ([]models.Contact, error) {
	var contacts []models.Contact
	query := `SELECT id, phone_number, email, linked_id, link_precedence, created_at, updated_at 
	          FROM contacts 
	          WHERE email = ? OR phone_number = ? OR linked_id IN 
	              (SELECT id FROM contacts WHERE email = ? OR phone_number = ?)`

	rows, err := database.DB.Query(query, email, phoneNumber, email, phoneNumber)
	if err != nil {
		log.Println("DB Query Error:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Contact
		if err := rows.Scan(&c.ID, &c.PhoneNumber, &c.Email, &c.LinkedID, &c.LinkPrecedence, &c.CreatedAt, &c.UpdatedAt); err != nil {
			log.Println("Row Scan Error:", err)
			continue
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func CreateContact(email, phoneNumber *string, linkedID *int, precedence string) int {
	query := `INSERT INTO contacts (phone_number, email, linked_id, link_precedence, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, NOW(), NOW())`

	// Execute the query
	result, err := database.DB.Exec(query, phoneNumber, email, linkedID, precedence)
	if err != nil {
		log.Printf("Error inserting contact (email: %v, phone: %v, linkedID: %v): %v", email, phoneNumber, linkedID, err)
		return 0
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return 0
	}

	return int(id)
}

func SegregateContacts(contacts []models.Contact) (models.Contact, []models.Contact) {
	if len(contacts) == 0 {
		return models.Contact{}, nil
	}

	primary := contacts[0]
	var secondaries []models.Contact

	for _, c := range contacts {
		if c.LinkPrecedence == "primary" && c.CreatedAt.Before(primary.CreatedAt) {
			primary = c
		}
	}

	for _, c := range contacts {
		if c.ID != primary.ID {
			secondaries = append(secondaries, c)
		}
	}

	return primary, secondaries
}

func FormatResponse(contacts []models.Contact) models.IdentifyResponse {
	var response models.IdentifyResponse
	primary, secondaries := SegregateContacts(contacts)

	response.Contact.PrimaryContactID = primary.ID

	emailSet := make(map[string]bool)
	phoneSet := make(map[string]bool)

	if primary.Email != nil {
		emailSet[*primary.Email] = true
	}
	if primary.PhoneNumber != nil {
		phoneSet[*primary.PhoneNumber] = true
	}

	for _, sec := range secondaries {
		response.Contact.SecondaryContactIDs = append(response.Contact.SecondaryContactIDs, sec.ID)
		if sec.Email != nil {
			emailSet[*sec.Email] = true
		}
		if sec.PhoneNumber != nil {
			phoneSet[*sec.PhoneNumber] = true
		}
	}

	for email := range emailSet {
		response.Contact.Emails = append(response.Contact.Emails, email)
	}
	for phone := range phoneSet {
		response.Contact.PhoneNumbers = append(response.Contact.PhoneNumbers, phone)
	}

	return response
}
