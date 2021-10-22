package data

type Users struct {
	ID           int
	Name         string `json:"name"`
	release_time int
	UpdateAt     string `json:"updateAt" sql:"not null;type:date"`
}

type Trains struct {
	TRAIN1 string
	TIME1  string
	TRAIN2 string
	TIME2  string
	TRAIN3 string
	TIME3  string
}

type Stations struct {
	Name           string
	First_Station  string
	Second_Station string
}
