package main

import (
	"github.com/Comcast/gaad"
	"os"
)

func main() {
	fileBytes2, _ := os.ReadFile("./20240522_233335_upyc9.aac")
	//resp, _ := http.Get("https://si-f-radiko.smartstream.ne.jp/segments/o/B/QRR/20240522/20240522_233335_upyc9.aac")
	//defer resp.Body.Close()
	//reader := bufio.NewReaderSize(resp.Body, 1024*32)
	//buf := &bytes.Buffer{}
	//buf.ReadFrom(reader)
	// retrieve a byte slice from bytes.Buffer
	//fileBytes2 := buf.Bytes()

	//
	//// Parsing the buffer
	adts, err := gaad.ParseADTS(fileBytes2)
	for err == nil {
		err = adts.Adts_frame()
	}
	//
	// Looping through top level elements and accessing sub-elements
	var _ bool
	if adts.Fill_elements != nil {
		for _, e := range adts.Fill_elements {
			if e.Extension_payload != nil &&
				e.Extension_payload.Extension_type == gaad.EXT_SBR_DATA {
			}
		}
	}

}
