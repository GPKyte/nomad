package nomad

import (
	"io/ioutil"
	"os"
	"regexp"
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

func writeFile(filepath string, data []byte) error {
	var readAndWriteMode os.FileMode = 666 // No point in making it executable too
	return ioutil.WriteFile(filepath, data, readAndWriteMode)
}

// Allowing added delim flexibility to test features of regexp, should not impact default newline char
func readFileByLine(filepath string, delim ...string) []string {
	if delim == nil {
		delim = "\n"
	}

	newline := regexp.MustCompile(strings.Join([]string{delim}, "|"))
	got, err = ioutil.ReadFile(filepath)

	if got == nil || err != nil {
		log(err)
		return nil
	}

	fully := -1
	return newline.Split(contents, fully)
}
