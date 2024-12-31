package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
)

func (s *Service) BackupDatabase(dbName string) error {
	godotenv.Load()
	dbHost, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	timestamp := time.Now().Format("20060102_150405")
	backupFile := fmt.Sprintf("./backup/%s_%s", dbName, timestamp)
	pwdOption := fmt.Sprintf("PGPASSWORD=%s", dbPwd)
	cmd := exec.Command("pg_dump",
		"-U", "postgres",
		"-h", dbHost,
		"-p", dbPort,
		"-F", "c",
		"-f", backupFile,
		dbName)
	cmd.Env = append(os.Environ(), pwdOption)
	err := cmd.Run()
	if err != nil {
		log.Printf("error executing pg_dump: %v", err)
		return err
	}
	log.Printf("backup completed successfully. file: %s", backupFile)
	return nil
}

func (s *Service) CronStart() {
	s.mycron.StartCron()
}

func (s *Service) CronStop() {
	s.mycron.StopCron()
}

func (s *Service) CronBackupDataBase(dbName string, query []string) error {
	err := s.mycron.AddJob(dbName, query, s)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveCronJob(id string) {
	s.mycron.RemoveJob(id)
}

func (s *Service) GetJobList() map[string]string {
	return s.mycron.GetJobList()
}
