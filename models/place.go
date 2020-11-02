package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Places .
type Places struct {
	Places []Place `json:"places"` 
}

//Place .
type Place struct{
	ID           primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Title        string                `json:"title"`
	Description  string                `json:"description"`
	Image        string                `json:"image"`
	Address      string                `json:"address"`
  Location     Location              `json:"location"`
	Creator      primitive.ObjectID    `json:"creator"`
}

//Location .
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `josn:"lng"`
}