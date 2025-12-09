package model

// ProgramResponse 节目 API 响应
type ProgramResponse struct {
	Stations []StationProgram `json:"stations"`
}

// StationProgram 电台节目信息
type StationProgram struct {
	StationID string   `json:"station_id"`
	Programs  Programs `json:"programs"`
}

// Programs 节目列表容器
type Programs struct {
	Date    string    `json:"date"`
	Program []Program `json:"program"`
}

// Program 单个节目信息
type Program struct {
	Ft    string `json:"ft"`    // 开始时间 YYYYMMDDHHMMSS
	To    string `json:"to"`    // 结束时间 YYYYMMDDHHMMSS
	Title string `json:"title"` // 节目标题
	Pfm   string `json:"pfm"`   // 主持人
}
