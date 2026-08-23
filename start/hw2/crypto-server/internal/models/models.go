package models

import (
	"time"
)

type User struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashedpassword"`
}

type Crypto struct {
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	CurrentPrice float64   `json:"currentprice"`
	LastUpdated  time.Time `json:"lastupdated"`
}

type PriceRecord struct {
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

type ScheduleConfig struct {
	Enabled         bool      `json:"enabled"`
	IntervalSeconds int       `json:"intervalseconds"`
	LastUpdate      time.Time `json:"lastupdate"`
	NextUpdate      time.Time `json:"nextupdate"`
}
