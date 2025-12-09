package model

// ProgramResponse represents the program API response
type ProgramResponse struct {
	Stations []StationProgram `json:"stations"`
}

// StationProgram represents station program information
type StationProgram struct {
	StationID string   `json:"station_id"`
	Programs  Programs `json:"programs"`
}

// Programs represents the program list container
type Programs struct {
	Date    string    `json:"date"`
	Program []Program `json:"program"`
}

// Program represents a single program
type Program struct {
	Ft    string `json:"ft"`    // Start time YYYYMMDDHHMMSS
	To    string `json:"to"`    // End time YYYYMMDDHHMMSS
	Title string `json:"title"` // Program title
	Pfm   string `json:"pfm"`   // Host/Performer
}
