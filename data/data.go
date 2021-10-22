package data

type Users struct {
	ID           int
	Name         string `json:"name"`
	release_time int
	UpdateAt     string `json:"updateAt" sql:"not null;type:date"`
}
