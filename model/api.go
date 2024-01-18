package model

type API struct {
	Path     string `json:"path,omitempty" bson:"path,omitempty"`
	IsActive *bool  `json:"is_active,omitempty" bson:"is_active,omitempty"`
}
