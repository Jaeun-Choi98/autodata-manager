package main

import (
	"cju/dao"
	"cju/entity"
	"log"
)

func main() {

	hobbies := []string{"Reading", "Cycling", "Gaming"}
	jobs := []entity.Job{
		{ID: 1, Name: "Software Engineer", UserId: 1},
		{ID: 2, Name: "Freelance Writer", UserId: 1},
	}
	user := entity.User{
		ID:      1,
		Age:     30,
		Name:    "John Doe",
		Hobbies: hobbies,
		Jobs:    jobs,
	}

	var mydb dao.DBLayerInterface
	con := "user=postgres dbname=test password=cjswo123 host=localhost port=5432 sslmode=disable"
	mydb, _ = dao.NewPostgreSQL(con)
	defer mydb.ClosePostgreSQL()
	err := mydb.AutoMigrateJob()
	if err != nil {
		log.Println(err)
	}
	mydb.AutoMigrateUesr()
	log.Println(user)
	mydb.AddUser(user)
}
