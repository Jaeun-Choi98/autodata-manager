package service

import (
	"fmt"
	"log"
	"strconv"

	"github.com/robfig/cron/v3"
)

type CronManager struct {
	mycron  *cron.Cron
	jobList map[string]string
	jobId   map[string]cron.EntryID
}

func NewCronInstance() *CronManager {
	return &CronManager{cron.New(), make(map[string]string), make(map[string]cron.EntryID)}
}

func (c *CronManager) StartCron() {
	c.mycron.Start()
}

func (c *CronManager) StopCron() {
	c.mycron.Stop()
}

func (c *CronManager) AddJob(dbName string, query []string, service *Service) error {
	if len(query) != 5 {
		log.Println("invaild cron query")
		return fmt.Errorf("invalid cron query")
	}
	nq := fmt.Sprintf("%s %s %s %s %s", query[0], query[1], query[2], query[3], query[4])
	_, err := cron.ParseStandard(nq)
	if err != nil {
		log.Printf("invaild cron query: %v", err)
		return err
	}
	id, _ := c.mycron.AddFunc(nq, func() {
		err = service.BackupDatabase(dbName)
		if err != nil {
			log.Println("cron backup DB error")
		}
	})
	if err != nil {
		return err
	}
	log.Printf("cron successful: %s->%s\n", dbName, query)
	strId := strconv.Itoa(int(id))
	c.jobList[strId] = fmt.Sprintf("%s -> %s", dbName, query)
	c.jobId[strId] = id
	return nil
}

func (c *CronManager) RemoveJob(id string) {
	c.mycron.Remove(cron.EntryID(c.jobId[id]))
	delete(c.jobList, id)
	delete(c.jobId, id)
}

func (c *CronManager) GetJobList() map[string]string {
	if len(c.jobList) == 0 {
		return map[string]string{"jobs": "nothing"}
	}
	return c.jobList
}
