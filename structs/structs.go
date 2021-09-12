package structs

type Logs struct {
	Id           int
	Trash_can_id int
	Type         string
	Timestamp    string
}

type Capacities struct {
	Id               int
	Trash_can_id     int
	Plastic_capacity float32
	Metal_capacity   float32
	Glass_capacity   float32
	Timestamp        string
}

type Types struct {
	Id        int
	Type_name string
}

type Trashcans struct {
	Id       int
	Location string
}
