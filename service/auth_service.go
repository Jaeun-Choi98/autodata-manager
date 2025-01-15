package service

import (
	"cju/entity/auth"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) UpdateUserFromCSV(filePath string) error {
	// read csv file
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, _ := reader.Read()
	records, _ := reader.ReadAll()
	users := CSVToUsers(&headers, &records)
	if err := s.mydb.UpdateUser(users); err != nil {
		return err
	}
	return nil
}

func (s *Service) AddUserFromCSV(filePath string) error {

	// read csv file
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, _ := reader.Read()
	records, _ := reader.ReadAll()
	users := CSVToUsers(&headers, &records)
	if err := s.mydb.AddUser(users); err != nil {
		return err
	}
	return nil
}

func CSVToUsers(headers *[]string, records *[][]string) []*auth.User {
	var users []*auth.User
	roleMap := make(map[string]auth.Role)

	for _, record := range *records {
		user := &auth.User{}
		profile := auth.Profile{}
		var roles []auth.Role
		for i, header := range *headers {
			switch header {
			case "Username":
				user.Username = record[i]
			case "Email":
				user.Email = record[i]
			case "Password":
				user.Password = record[i]
			case "IsActive":
				isActive, _ := strconv.ParseBool(record[i])
				user.IsActive = isActive
			case "FirstName":
				profile.FirstName = record[i]
			case "LastName":
				profile.LastName = record[i]
			case "PhoneNumber":
				profile.PhoneNumber = record[i]
			case "Address":
				profile.Address = record[i]
			case "ProfilePicture":
				profile.ProfilePicture = record[i]
			case "Roles":
				roleNames := strings.Split(record[i], ",")
				for _, roleName := range roleNames {
					roleName := strings.TrimSpace(roleName)
					if _, exists := roleMap[roleName]; !exists {
						roleMap[roleName] = auth.Role{RoleName: roleName}
					}
					roles = append(roles, roleMap[roleName])
				}
			}
		}
		user.Profile = profile
		user.Roles = roles
		users = append(users, user)
	}
	return users
}

func hashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword를 사용해 해시 생성
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func checkPassword(hashedPassword, password string) bool {
	// bcrypt.CompareHashAndPassword로 검증
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
