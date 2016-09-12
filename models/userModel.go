package models

type User struct {
	User     int    `json:"User" bson:"User"`
	GoogleID string `json:"GoogleID" bson:"GoogleID"`
	Name     string `json:"Name" bson:"Name"`
}
