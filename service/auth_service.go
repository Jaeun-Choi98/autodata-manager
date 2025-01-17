package service

import (
	"cju/entity/auth"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) ReadUserByEmail(email string) (*auth.User, error) {
	user, err := s.mydb.ReadUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(email, pwd string) (string, error) {
	user, err := s.mydb.ReadUserByEmail(email)
	if err != nil {
		return "", err
	}
	if !checkPassword(user.Password, pwd) {
		log.Println("password does not match")
		return "", fmt.Errorf("password does not match")
	}
	jwtString, err := GenerateJWT(user.Email, user.Roles[0].RoleName)
	if err != nil {
		return "", err
	}
	return jwtString, nil
}

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
	users, err := CSVToUsers(&headers, &records)
	if err != nil {
		return err
	}
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
	users, err := CSVToUsers(&headers, &records)
	if err != nil {
		return err
	}
	if err := s.mydb.AddUser(users); err != nil {
		return err
	}
	return nil
}

func CSVToUsers(headers *[]string, records *[][]string) ([]*auth.User, error) {
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
				pwd, err := hashPassword(record[i])
				if err != nil {
					return nil, err
				}
				user.Password = pwd
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
	return users, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// registerd claim( exp, iss, sub, aud, etc.. )
// 이후 user의 role이 여러 개일 수도 있음.
func GenerateJWT(email, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  email,
		"exp":  time.Now().Add(time.Hour * 12).Unix(),
		"role": role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_KEY")
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("failed to sign token: %v", err)
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString string) (string, string, error) {
	secretKey := os.Getenv("JWT_KEY")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Printf("failed to parse token: %v", err)
		return "", "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, _ := claims.GetSubject()
		role := claims["role"].(string)
		return email, role, nil
	}
	log.Println("invalid token")
	return "", "", fmt.Errorf("invalid token")
}
