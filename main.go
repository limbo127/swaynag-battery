package main

import (
	"fmt"
	"os"
	"time"
)

func DesiredDisplays(displays StringSet, activeDisplays StringSet) StringSet {
	if len(displays) == 0 {
		return activeDisplays
	}

	return Intersection(displays, activeDisplays)
}

func tick(watcher *Watcher, params Parameters) {
	battery, err := LoadBatteryInfo(params.uevent)
	if err != nil {
		logWarning("Skipping this cycle due to errors occurred.")
		return
	}

	displays := DesiredDisplays(params.displays, ActiveDisplays())

	if !battery.Charging() && battery.Capacity <= params.threshold {
		messages := ShowAll(params.message+" ["+fmt.Sprintf("%v", battery.Capacity)+"/ 100 ]", watcher.MessagesFor(displays))
		watcher.Update(messages, battery.Status)
	}

	if battery.Charging() {
		messages := watcher.Messages()
		CloseAll(messages)
		watcher.Empty()
		watcher.CleanUp(displays)
	}
}

func main() {
	params := CommandLineParameters(os.Args[1:])
	watcher := NewWatcher()

	tick(&watcher, params)
	for range time.Tick(params.interval) {
		tick(&watcher, params)
	}
}
