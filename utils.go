package nomad

import (
	"io/ioutil"
	"strings"
)

func loadAllegiantLocationData(index int) []location {
	var allegiantLocationData []byte = ioutil.ReadFile("resources/allegiant/locations.txt")

	newlineMarkers := make(chan int)
	lines := make(chan string)
	done := make(chan bool)

	// Parse each location into data struct on init
	go func() {
		locations := make([]string, 100)

		for {
			locData, more := <-lines

			if more {
				commaIndex := strings.IndexByte(locData, byte(','))
				afterParenIndex := strings.IndexByte(locData, byte('(')) + 1 // Don't want paren in the code

				name := locData[:commaIndex]
				state := locData[commaIndex+len(", ") : commaIndex+len(", XY")] // Get 'XY'
				code := locData[afterParenIndex : afterParenIndex+len("ABC")]   // Get 'ABC'

				locations = append(locations, location{Name: name, State: state, Code: code})
			} else {
				done <- true
				return
			}
		}
	}()

	go func() {
		start := 0
		for {
			nextDelim, more := <-newlineMarkers

			if more {
				lines <- string(allegiantLocationData[start:nextDelim])
				start = nextDelim + 1

			} else {
				lines <- string(allegiantLocationData[start:])
				return
			}
		}
	}()

	// Find all newline characters deliminating lines
	for i := 0; i < len(allegiantLocationData); i++ {
		if allegiantLocationData[i] == '\n' {
			newlineMarkers <- i
		}
	}

	<-done
	return locations
}
