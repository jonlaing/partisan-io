package cron

import (
	"fmt"
	"testing"
	"time"
)

func say(s string, count *int) func() {
	return func() {
		*count++
		fmt.Println(time.Now(), s)
	}
}

func lazySay(s string, count *int) func() {
	return func() {
		time.Sleep(5 * time.Second)
		*count++
		fmt.Println(time.Now(), s)
	}
}

func TestCron(t *testing.T) {
	var helloCount, whatsCount, lazyCount int

	sched := Scheduler{}
	sched.Seconds(1, say("hello", &helloCount))
	sched.Seconds(10, say("what's up", &whatsCount))
	sched.Seconds(10, lazySay("being lazy", &lazyCount))
	sched.Start()
	time.Sleep(25 * time.Second)
	sched.Stop()

	if helloCount != 25 {
		t.Error("expected 25 hellos, got:", helloCount)
	}

	if whatsCount != 2 {
		t.Error("Expected 2 what's ups, got:", whatsCount)
	}

	if lazyCount != 2 {
		t.Error("Expected 2 lazys, got:", lazyCount)
	}
}
