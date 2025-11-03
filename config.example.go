package main

import "fmt"

// 配置示例 - 可以根据需要修改这些值

// StationConfig 电台配置
type StationConfig struct {
	Name      string
	StationID string
}

// 预定义的电台列表
var Stations = map[string]StationConfig{
	"QRR": {
		Name:      "文化放送",
		StationID: "QRR",
	},
	"TBS": {
		Name:      "TBS ラジオ",
		StationID: "TBS",
	},
	"LFR": {
		Name:      "ニッポン放送",
		StationID: "LFR",
	},
	"INT": {
		Name:      "interfm",
		StationID: "INT",
	},
	"FMT": {
		Name:      "TOKYO FM",
		StationID: "FMT",
	},
	"FMJ": {
		Name:      "J-WAVE",
		StationID: "FMJ",
	},
}

// GetStationURL 根据电台 ID 生成完整 URL
func GetStationURL(stationID string, lsid string) string {
	baseURL := "https://c-radiko.smartstream.ne.jp/%s/_definst_/simul-stream.stream/playlist.m3u8"
	return fmt.Sprintf(baseURL+"?station_id=%s&l=30&lsid=%s&type=b", stationID, stationID, lsid)
}
