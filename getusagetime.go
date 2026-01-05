package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	shellCommand := `pmset -g log | grep "Display is turned"`
	cmd := exec.Command("bash", "-c", shellCommand)
	o, _ := cmd.Output()
	// output is []string
	output := strings.Split(string(o), "\n")
	// give the AC unplug date and time manually
	// yeah, but works for me.
	//txt := "2025-12-25 22:13:00 +0300"
	shellCommand2 := `pmset -g log | awk '/Using AC/ {last_was_ac=1} /Using (Batt|Battery)/ {if(last_was_ac) {print $0; last_was_ac=0}}' | tail -n 1`
	cmd2 := exec.Command("bash", "-c", shellCommand2)
	text_string, _ := cmd2.Output()
	txt := string(text_string)
	// this date is spesifically selected by golang developers.
	// used for date template
	layout := "2006-01-02 15:04:05 -0700"
	// a time we got from last charged date log
	beginning, _ := time.Parse(layout, txt[:25])
	// this will keep if screen is on
	var onScreen = beginning
	var totalTime float64
	for _, line := range output {
		if line == "" {
			continue
		}
		dateLine := strings.TrimSpace(line)
		curDate, _ := time.Parse(layout, dateLine[:25])
		// if date earlier than last charge
		// discard
		if curDate.Before(beginning) {
			continue
		}
		if strings.Contains(dateLine, " on") {
			onScreen = curDate
		} else {
			if !onScreen.IsZero() {
				totalTime += curDate.Sub(onScreen).Hours()
			}
			onScreen = time.Time{}
		}
	}
	if !onScreen.IsZero() {
		totalTime += time.Now().Sub(onScreen).Hours()
	}
	fmt.Println("Total screen on battery usage:", totalTime, " hours")
}
