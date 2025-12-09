package model

import "encoding/xml"

type RadikoStations struct {
	XMLName  xml.Name  `xml:"radiko"`
	Stations []Station `xml:"stations>station"`
}

type Station struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

type RadikoURLs struct {
	XMLName xml.Name `xml:"urls"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
	AreaFree          int    `xml:"areafree,attr"`
	TimeFree          int    `xml:"timefree,attr"`
	PlaylistCreateURL string `xml:"playlist_create_url"`
}
