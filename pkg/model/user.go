package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Action int

const (
	// Actions
	Buy Action = iota
	Sell
)

type ExperienceLevel string

const (
	// Experience Levels
	Beginner ExperienceLevel = "Beginner"
	Advanced ExperienceLevel = "Advanced"
)

type Allocation string

const (
	ShortTerm ExperienceLevel = "ShortTerm"
	LongTerm  ExperienceLevel = "LongTerm"
)

type RiskLevel string

const (
	Low  RiskLevel = "Low"
	High RiskLevel = "High"
)

type User struct {
	gorm.Model
	PhoneNumber string `gorm:"uniqueIndex"`
	Email       string `gorm:"uniqueIndex"`
	FirebaseUID string `gorm:"uniqueIndex"`
	Username    string
	Experience  ExperienceLevel
	Allocation  Allocation
	Risk        RiskLevel
	Onboarded   bool
}

type Watchlist struct {
	gorm.Model
	UserID  uint
	Name    string
	User    User           `gorm:"foreignKey:UserID"`
	Symbols pq.StringArray `gorm:"type:text[]"`
}

type Notification struct {
	gorm.Model
	UserID  uint
	Action  Action
	Message string
	User    User `gorm:"foreignKey:UserID"`
}

type ChatMessage struct {
	gorm.Model
	Question string
	Answer   string
	UserID   uint
	User     User `gorm:"foreignKey:UserID"`
}

type Stragegy struct{}

type Stock struct{}
