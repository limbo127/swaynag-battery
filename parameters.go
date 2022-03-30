package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	docopt "github.com/docopt/docopt-go"
)

type Parameters struct {
	displays  StringSet
	interval  time.Duration
	message   string
	threshold int
	uevent    string
}

var (
	usage = `
Shows a message (using swaynag) when battery percentage is less then specified
value.

Usage:
  swaynag-battery [options]
  swaynag-battery -h | --help
  swaynag-battery --version

Options:
  --displays <display-list>  Comma separated list of displays to show the
                             alert - the default is to show in all displays.
  --threshold <int>          Percentual threshold to show notification.
                             [default: 15]
  --message <string>         Message to display [default: You battery is running low. Please plug in a power adapter] 
  --interval <duration>      Check battery at every interval. [default: 5m]
  --uevent <path>            Uevent path for reading battery stats.
                             [default: auto]
  -h --help                  Show this screen.
  --version                  Show version.

`
)

func isBattery(path string) bool {
	t, err := ioutil.ReadFile(filepath.Join(path, "type"))
	return err == nil && string(t) == "Battery\n"
}

func findBattery() string {
	const sysfs = "/sys/class/power_supply"
	files, err := ioutil.ReadDir(sysfs)
	if err != nil {
		return ""
	}
	for _, file := range files {
		match, _ := regexp.MatchString("hid.*", file.Name())
		path := filepath.Join(sysfs, file.Name())
		if !match && isBattery(path) {
			return path + "/uevent"
		}
	}
	return ""
}

func CommandLineParameters(arguments []string) Parameters {
	args, err := docopt.ParseArgs(usage, arguments, version)
	if err != nil {
		logAndExit(18, "Unable to parse input arguments.")
	}

	interval, err := time.ParseDuration(args["--interval"].(string))
	if err != nil {
		logAndExit(28, "Unable to parse '--interval %s': the value must be a duration.", args["--interval"])
	}

	threshold, err := strconv.Atoi(args["--threshold"].(string))
	if err != nil {
		logAndExit(38, "Unable to parse '--threshold %s': the value must be an integer number.", args["--threshold"])
	}

	displays := []string{}
	d, ok := args["--displays"].(string)
	if ok {
		displays = strings.Split(d, ",")
	}

	uevent := args["--uevent"].(string)
	if uevent == "auto" {
		uevent = findBattery()
	}
	file, err := os.Open(uevent)
	if err != nil {
		logAndExit(42, "Could not load battery file '%s'.", uevent)
	}
	file.Close()

	message := args["--message"].(string)

	return Parameters{
		displays:  SetFrom(displays),
		interval:  interval,
		message:   message,
		threshold: threshold,
		uevent:    uevent}
}
