package auth

import (
	"time"
)

type User struct {
	ID        uint        `gorm:"primaryKey;autoIncrement"`
	Username  string      `gorm:"type:varchar(50);unique;not null"`
	Email     string      `gorm:"type:varchar(100);unique;not null"`
	Password  string      `gorm:"type:text;not null"`
	IsActive  bool        `gorm:"default:true"`
	CreatedAt time.Time   `gorm:"default:current_timestamp"`
	UpdatedAt time.Time   `gorm:"default:current_timestamp"`
	Profile   UserProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Roles     []UserRole  `gorm:"many2many:user_role_assignments;constraint:OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "auth.users"
}

type UserProfile struct {
	UserID         uint   `gorm:"primaryKey"`
	FirstName      string `gorm:"type:varchar(50)"`
	LastName       string `gorm:"type:varchar(50)"`
	PhoneNumber    string `gorm:"type:varchar(15)"`
	Address        string `gorm:"type:text"`
	ProfilePicture string `gorm:"type:text"`
}

func (UserProfile) TableName() string {
	return "auth.user_profiles"
}

type UserRole struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	RoleName    string `gorm:"type:varchar(50);unique;not null"`
	Description string `gorm:"type:text"`
	Users       []User `gorm:"many2many:user_role_assignments;constraint:OnDelete:CASCADE"`
}

func (UserRole) TableName() string {
	return "auth.user_roles"
}

// Many-to-Many
type UserRoleAssignment struct {
	UserID     uint      `gorm:"primaryKey"`
	RoleID     uint      `gorm:"primaryKey"`
	AssignedAt time.Time `gorm:"default:current_timestamp"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Role       UserRole  `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}

func (UserRoleAssignment) TableName() string {
	return "auth.user_role_assignments"
}
