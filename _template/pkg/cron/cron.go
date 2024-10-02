package cron

import (
	"github.com/robfig/cron"
)

// type cronMethod with any params and output
type CronMethod func()

type CronJob struct {
	interval string
	method   CronMethod
}

func NewCronJob(interval string, method CronMethod) *CronJob {
	return &CronJob{
		interval: interval,
		method:   method,
	}
}

func Start(jobList []CronJob) {
	c := cron.New()

	AddJobList(c, jobList)

	c.Start()
}

func AddJobList(c *cron.Cron, jobList []CronJob) {

	for _, job := range jobList {
		c.AddFunc(job.interval, job.method)
	}

}
