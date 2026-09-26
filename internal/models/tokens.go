package models

import "time"

type Tokens struct {
	Access   string
	Refresh  string
	JTIAcc   string
	JTIRef   string
	ExpAcc   time.Time
	ExpRef   time.Time
	UserID   string
	Issuer   string
	Audience string
}
