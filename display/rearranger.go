package display

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

const (
	scanEvery = 2 * time.Second
)

var displayCount int
var preferences []Preference

func init() {
	displayCount = 0
	readPreferences()
}

func Start() {
	scheduler := Scheduler{}
	scheduler.Do(checkDisplays).Every(scanEvery).Run()
}

func checkDisplays() {
	log.Println("checkDisplays - Running...")
	rearrange()
}

func currentOutputs() []string {
	return scanConnectedOutputs()
}

func rearrange() {
	outputs := currentOutputs()
	log.Println("rearrange - Outputs:", outputs)

	newCount := len(outputs)
	if newCount != displayCount {
		log.Println("rearrange - newCount:,", newCount, " != displayCount:", displayCount)

		if arrangement, found := findArrangement(outputs); found {
			log.Println("rearrange - arrangement:", arrangement)
			if applyArrangement(arrangement) {
				displayCount = newCount
			}
		}
	}
}

func findArrangement(outputs []string) (Arrangement, bool) {
	for _, preference := range preferences {
		if compareSlice(outputs, preference.When) {
			return preference.Arrangement, true
		}
	}
	return Arrangement{}, false
}

func compareSlice(a []string, b []string) bool {
	fmt.Println("compareSlice:", a, b)
	return reflect.DeepEqual(a, b)
}

func applyArrangement(arrangement Arrangement) bool {
	fmt.Println("arrangement", arrangement)
	for _, display := range arrangement {
		if err := arrangeDisplay(display); err != nil {
			fmt.Println("Cannot arrange display", display, "due:", err)
			return false
		}
	}
	return true
}

func makeXrandrScanCommand() string {
	return fmt.Sprintf("xrandr | grep ' connected' | awk '{print $1}'")
}

func makeXrandrApplyCommand(display Display) string {
	return fmt.Sprintf("xrandr "+
		" --output %s"+
		" --mode %dx%d"+
		" --pos %dx%d"+
		" --scale %dx%d"+
		" --auto",
		display.Output,
		display.Mode.Width, display.Mode.Height,
		display.Pos.X, display.Pos.Y,
		display.Scale, display.Scale,
	)
}

func arrangeDisplay(display Display) error {
	cmd := makeXrandrApplyCommand(display)
	out, err := exec.Command("sh", "-c", cmd).Output()
	fmt.Println("arrangeDisplay - cmd: ", cmd)
	fmt.Println("arrangeDisplay - out: ", string(out))
	return err
}

func scanConnectedOutputs() (outputs []string) {
	cmd := makeXrandrScanCommand()
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Println("Cannot read connected displays:", err)
		return
	}
	lines := strings.Trim(string(out), "\n")
	outputs = strings.Split(lines, "\n")
	return
}

func readPreferences() {
	home := os.Getenv("HOME")
	path := filepath.Join(home, "display-preferences.json")
	log.Println("Opening config file:", path)
	f, err := os.Open(path)
	if err != nil {
		log.Panic("Error while loading config file.", err)
	}
	defer f.Close()
	content, _ := ioutil.ReadAll(f)
	if err := json.Unmarshal(content, &preferences); err != nil {
		log.Panic("Cannot understand the config file.", err)
	}
}
