package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	timeUntilSleepAsString := os.Args[1]

	timeParts := strings.Split(timeUntilSleepAsString, ":")
	timePartsLen := len(timeParts)

	if timePartsLen > 2 || !checkFormat(timeUntilSleepAsString) {
		err := fmt.Errorf("expected format hh:mm, got %s", timeUntilSleepAsString)
		printErr(err)
	}

	// By default just minutes are assumed
	minutes, err := strconv.Atoi(timeParts[0])
	if err != nil {
		printErr(err)
	}

	if timePartsLen == 2 {
		// If the string contains hours and minutes,
		// convert the hours to minutes
		minutes = minutes * 60
		// Add the minutes part
		convertedMinutes, err := strconv.Atoi(timeParts[1])
		if err != nil {
			printErr(err)
		}

		minutes += convertedMinutes
	}

	minuteMsg := "minutes"
	if minutes == 1 {
		minuteMsg = "minute"
	}

	fmt.Printf("Going to sleep in %d %s\n", minutes, minuteMsg)
	time.Sleep(time.Duration(minutes) * time.Minute)
	err = goToSleepDarwin()
	if err != nil {
		printErr(err)
	}
}

func goToSleepDarwin() error {
	// Kinda sucks because the App has to be allowed to send messages to System Events
	// sleepCmd := exec.Command("osascript", "-e", `tell app "System Events" to sleep`)
	// Therefore pmset is used
	sleepCmd := exec.Command("pmset", "sleepnow")

	err := sleepCmd.Start()
	if err != nil {
		return err
	}

	err = sleepCmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func checkFormat(input string) bool {
	formatRegex := regexp.MustCompilePOSIX(`(([0123456789]{1,2}:)?[0123456789]{1,2})`)
	// Use find instead of match, because match detects substrings as well
	matchedInput := formatRegex.FindString(input)
	return input == matchedInput
}

func printUsage() {
	printErr(nil)
}

func printErr(err error) {
	fmt.Printf("Usage: %s [hh:]mm\n", filepath.Base(os.Args[0]))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(1)
}
