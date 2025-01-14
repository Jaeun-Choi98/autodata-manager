package auth

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"type:varchar(50);unique;not null"`
	Email     string    `gorm:"type:varchar(100);unique;not null"`
	Password  string    `gorm:"type:text;not null"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
	Profile   Profile   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Roles     []Role    `gorm:"many2many:auth.user_roles;constraint:OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "auth.users"
}

type Profile struct {
	UserID         uint   `gorm:"primaryKey"`
	FirstName      string `gorm:"type:varchar(50)"`
	LastName       string `gorm:"type:varchar(50)"`
	PhoneNumber    string `gorm:"type:varchar(20)"`
	Address        string `gorm:"type:text"`
	ProfilePicture string `gorm:"type:text"`
}

func (Profile) TableName() string {
	return "auth.profiles"
}

type Role struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	RoleName    string `gorm:"type:varchar(50);unique;not null"`
	Description string `gorm:"type:text"`
	Users       []User `gorm:"many2many:auth.user_roles;constraint:OnDelete:CASCADE"`
}

func (Role) TableName() string {
	return "auth.roles"
}

// Many-to-Many
/*
	gorm many2many 태그를 사용하지 않고, 아래처럼 사용자 정의해서 사용할 수도 있음.
	단, User, Role에 many2many 태그는 사용x
type UserRole struct {
	UserID     uint      `gorm:"primaryKey"`
	RoleID     uint      `gorm:"primaryKey"`
	AssignedAt time.Time `gorm:"default:current_timestamp"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Role       Role      `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}

func (UserRole) TableName() string {
	return "auth.user_roles"
}
*/
