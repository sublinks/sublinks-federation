package model

type Actor struct {
	ActorType    string `json:"actor_type" gorm:"primarykey"`
	Id           string `json:"id" gorm:"primarykey"`
	Username     string `json:"username,omitempty" gorm:"not null"`
	Name         string `json:"name,omitempty" gorm:"nullable"`
	Bio          string `json:"bio"`
	MatrixUserId string `json:"matrix_user_id,omitempty"`
	PublicKey    string `json:"public_key"`
	PrivateKey   string `json:"private_key"`
}
