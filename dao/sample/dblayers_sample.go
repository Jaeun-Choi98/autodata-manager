package dao

import entity "cju/entity/sample"

type DBLayerInterface interface {
	AutoMigrateJob() error
	AutoMigrateUesr() error
	AddUser(user entity.User) error
	GetUser() (*entity.User, error)
	ClosePostgreSQL() error
}
