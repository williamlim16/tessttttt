package structs

import "time"

// Structs according to Database (Snake_case)

type User struct {
	Id           int
	Name         string
	Email        string
	Telephone    string
	Address      string
	Company_name string
	Password     string
	Created_at   time.Time
	Updated_at   time.Time
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
	Created_at             time.Time
	Assigned_date          time.Time
	User_id                int
	Latitude               string
	Longitude              string
	Location               string
	Guarantee_expired_date string
	Trash_version_id       int
}

type Trash_capacity struct {
	Id                 int
	Trash_id           int
	Organic_capacity   int
	Anorganic_capacity int
	Created_at         time.Time
}

type Trash_reading struct {
	Id         int
	Trash_id   int
	Category   string
	Type       string
	Created_at time.Time
}

// Structs for responses (camelCase)

type TrashCapacity struct {
	Trash_can_id         int
	Organic_capacity     int
	Anorganic_capacity   int
	Organic_max_height   int
	Anorganic_max_height int
}
