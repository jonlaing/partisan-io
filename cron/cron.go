package cron

import (
	"sort"
	"time"
)

// Units to figure out how much time to add. I do this instead of
// `time.Interval`, because months don't have a consistent interval,
// so I take advantage of time.AddDate() for that, but it doesn't
// take seconds, minutes or hours as parameters, so I still have to
// use intervals for those.
const (
	seconds uint = iota
	minutes
	hours
	days
	weeks
	months
	years
)

type job struct {
	task     func()
	interval int
	unit     uint
	nextRun  time.Time
}

func newJob(task func(), i int, unit uint) *job {
	job := job{
		task:     task,
		interval: i,
		unit:     unit,
	}

	job.scheduleNext()
	return &job
}

func (j *job) scheduleNext() {
	if j.nextRun.IsZero() {
		j.nextRun = time.Now()
	}

	interval := time.Duration(j.interval)

	switch j.unit {
	case minutes:
		j.nextRun = j.nextRun.Add(interval * time.Minute)
	case hours:
		j.nextRun = j.nextRun.Add(interval * time.Hour)
	case days:
		j.nextRun = j.nextRun.AddDate(0, 0, j.interval)
	case weeks:
		j.nextRun = j.nextRun.AddDate(0, 0, j.interval*7)
	case months:
		j.nextRun = j.nextRun.AddDate(0, j.interval, 0)
	case years:
		j.nextRun = j.nextRun.AddDate(j.interval, 0, 0)
	default:
		j.nextRun = j.nextRun.Add(interval * time.Second)
	}
}

func (j job) shouldRun() bool {
	return time.Now().After(j.nextRun)
}

// Scheduler holds all of the jobs and a ticker, as well as a "stop" channel.
type Scheduler struct {
	jobs   []*job
	ticker *time.Ticker
	stop   chan bool
}

// Len implements the Sorter interfaces
func (s *Scheduler) Len() int {
	return len(s.jobs)
}

// Swap implements the Sorter interface
func (s *Scheduler) Swap(i, j int) {
	s.jobs[i], s.jobs[j] = s.jobs[j], s.jobs[i]
}

// Less implements the Sorter interface. We order it on the assumption that tasks that
// are meant to be run more frequently take less time, and have more imperative to run
// exactly when they are scheduled than tasks that are scheduled to happen less frequently.
func (s *Scheduler) Less(i, j int) bool {
	if s.jobs[i].unit == s.jobs[j].unit {
		return s.jobs[i].interval < s.jobs[j].interval
	}

	return s.jobs[i].unit < s.jobs[j].unit
}

func (s *Scheduler) addJob(task func(), i int, unit uint) {
	s.jobs = append(s.jobs, newJob(task, i, unit))
	sort.Sort(s)
}

// Seconds allows you to set a task to run every specified number of seconds.
func (s *Scheduler) Seconds(i int, task func()) {
	s.addJob(task, i, seconds)
}

// Minutes allows you to set a task to run every  specified number of seconds.
func (s *Scheduler) Minutes(i int, task func()) {
	s.addJob(task, i, minutes)
}

// Hours allows you to set a task to run every  specified number of seconds.
func (s *Scheduler) Hours(i int, task func()) {
	s.addJob(task, i, hours)
}

// Days allows you to set a task to run every  specified number of seconds.
func (s *Scheduler) Days(i int, task func()) {
	s.addJob(task, i, days)
}

// Days allows you to set a task to run every  specified number of seconds.
func (s *Scheduler) Weeks(i int, task func()) {
	s.addJob(task, i, weeks)
}

func (s *Scheduler) Months(i int, task func()) {
	s.addJob(task, i, months)
}

func (s *Scheduler) Years(i int, task func()) {
	s.addJob(task, i, years)
}

func (s *Scheduler) Start() {
	s.stop = make(chan bool, 1)
	s.ticker = time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.runPending()
			case <-s.stop:
				return
			}
		}
	}()
}

func (s *Scheduler) Clear() {
	s.jobs = []*job{}
}

func (s *Scheduler) Stop() {
	s.ticker.Stop()
	s.stop <- true
}

func (s *Scheduler) runPending() {
	for i := range s.jobs {
		if s.jobs[i].shouldRun() {
			go s.jobs[i].task()
			s.jobs[i].scheduleNext()
		}
	}
}
