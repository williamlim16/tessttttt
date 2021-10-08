package structs

import "time"

// Structs according to Database (Snake_case)

type User struct {
	Id           int
	Name         string `json:"name"`
	Email        string `json:"email"`
	Telephone    string
	Address      string
	Company_name string
	Password     string `json:"password"`
	ConfirmPass  string `json:"confirm_password"`
	Created_date time.Time
	Updated_date time.Time
}

type Trash_version struct {
	Id                   int
	Version_name         string
	Inorganic_max_height int
	Organic_max_height   int
}

type Trash struct {
	Id                     int
	Trash_code             string
	Assigned               string
	Created_date           time.Time
	Assigned_date          time.Time
	User_id                int
	Latitude               string
	Longitude              string
	Location               string
	Guarantee_expired_date string
	Trash_version_id       int
	Custom_name            string
}

type Trash_capacity struct {
	Id                 int
	Trash_id           int
	Organic_capacity   int
	Inorganic_capacity int
	Created_at         time.Time
}

type Trash_reading struct {
	Id           int
	Trash_id     int
	Category     string
	Type         string
	Created_date time.Time
}

// Structs for responses (camelCase)

type TrashCapacity struct {
	Trash_id             int
	Organic_capacity     int
	Inorganic_capacity   int
	Organic_max_height   int
	Inorganic_max_height int
}

type TrashLogs struct {
	Trash_can_id  int
	Trash_reading []Trash_reading
}

type TrashReading struct {
	Trash_sorter_name     string
	Trash_sorter_location string
	Total                 int
}
