package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"radiko-tui/model"
)

const (
	StationListURLFmt = "https://api.radiko.jp/program/v3/now/%s.xml"
	StreamURLFmt      = "https://radiko.jp/v3/station/stream/pc_html5/%s.xml"
)

// GetStations retrieves the list of stations for a specified area
func GetStations(areaID string) ([]model.Station, error) {
	url := fmt.Sprintf(StationListURLFmt, areaID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch station list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch station list: status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var radikoStations model.RadikoStations
	if err := xml.Unmarshal(data, &radikoStations); err != nil {
		return nil, fmt.Errorf("failed to parse station list XML: %w", err)
	}

	return radikoStations.Stations, nil
}

func GetStreamURLs(stationID string) ([]string, error) {
	url := fmt.Sprintf(StreamURLFmt, stationID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stream URL for station %s: %w", stationID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch stream URL: status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var radikoURLs model.RadikoURLs
	if err := xml.Unmarshal(data, &radikoURLs); err != nil {
		return nil, fmt.Errorf("failed to parse stream URL XML: %w", err)
	}

	if len(radikoURLs.URLs) == 0 {
		return nil, fmt.Errorf("no stream URLs found for station %s", stationID)
	}

	var urls []string
	for _, u := range radikoURLs.URLs {
		if u.PlaylistCreateURL != "" {
			urls = append(urls, u.PlaylistCreateURL)
		}
	}

	return urls, nil
}

// ProgramURLFmt is the program info API URL format
const ProgramURLFmt = "https://api.radiko.jp/program/v4/date/%s/station/%s.json"

var jst *time.Location

func init() {
	// Use Japan timezone (UTC+9)
	jst = time.FixedZone("JST", 9*60*60)
}

// GetCurrentProgram retrieves the current program for a station
func GetCurrentProgram(stationID string) (*model.Program, error) {
	now := time.Now().In(jst)
	dateStr := now.Format("20060102")
	timeStr := now.Format("20060102150405")

	// Try to get program for current date
	prog, err := getProgramForDate(stationID, dateStr, timeStr)
	if err != nil {
		return nil, err
	}

	// If program found, return it
	if prog != nil {
		return prog, nil
	}

	// If no program found, the API might have returned next day's data
	// Try fetching previous day's data
	yesterday := now.AddDate(0, 0, -1)
	yesterdayStr := yesterday.Format("20060102")

	prog, err = getProgramForDate(stationID, yesterdayStr, timeStr)
	if err != nil {
		return nil, err
	}

	return prog, nil
}

// getProgramForDate retrieves program data for a specific date and finds the current program
func getProgramForDate(stationID, dateStr, timeStr string) (*model.Program, error) {
	url := fmt.Sprintf(ProgramURLFmt, dateStr, stationID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var progResp model.ProgramResponse
	if err := json.Unmarshal(data, &progResp); err != nil {
		return nil, err
	}

	// Find the program for the current time
	for _, station := range progResp.Stations {
		if station.StationID == stationID {
			// Check if the first program starts after current time
			// This indicates the current program data is in the previous day's API data
			if len(station.Programs.Program) > 0 {
				firstProgram := station.Programs.Program[0]
				if firstProgram.Ft > timeStr {
					// First program is in the future, current program data is from previous day
					return nil, nil
				}
			}

			// Find the program that matches current time
			for _, prog := range station.Programs.Program {
				// Check if current time is within the program's time range
				if prog.Ft <= timeStr && timeStr < prog.To {
					return &prog, nil
				}
			}
		}
	}

	return nil, nil
}
