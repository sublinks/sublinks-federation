package model

type Actor struct {
	ActorType    string `json:"actor_type" gorm:"index"`
	Id           string `json:"id" gorm:"primarykey"`
	Username     string `json:"username"`
	Name         string `json:"name,omitempty" gorm:"nullable"`
	Bio          string `json:"bio"`
	MatrixUserId string `json:"matrix_user_id,omitempty"`
	PublicKey    string `json:"public_key"`
	PrivateKey   string `json:"private_key"`
}
