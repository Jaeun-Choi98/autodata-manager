package dao

import (
	"cju/entity/auth"
	"cju/utils"
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQL struct {
	db *gorm.DB
}

func NewPostgreSQL(con string) (*PostgreSQL, error) {
	db, err := gorm.Open(postgres.Open(con), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println("failed new db")
		return nil, err
	}
	return &PostgreSQL{db}, err
}

func (pg *PostgreSQL) CloseDB() error {
	db, err := pg.db.DB()
	if err != nil {
		log.Println("failed to close db")
		return err
	}
	return db.Close()
}

func (pq *PostgreSQL) ExecQuery(query string) error {
	err := pq.db.Exec(query).Error
	if err != nil {
		log.Println("failed to Exec Query")
		return err
	}
	return nil
}

func (pq *PostgreSQL) Init() error {
	var roles []*auth.Role
	if err := pq.db.Raw(`SELECT * FROM AUTH.ROLES`).Scan(&roles).Error; err != nil {
		log.Println(err)
		return err
	}
	utils.UploadRoleId(roles)
	return nil
}

func (pq *PostgreSQL) ExistTable(tableName string) bool {
	exists := pq.db.Migrator().HasTable(tableName)
	return exists
}

func (pq *PostgreSQL) ReadAllTableData(tableName string) ([]map[string]interface{}, error) {
	var rows []map[string]interface{}
	err := pq.db.Table(tableName).Find(&rows).Error
	if err != nil {
		log.Printf("failed to query table data: %v", err)
		return nil, err
	}
	return rows, nil
}

func (pq *PostgreSQL) ReadAllTables(schemaName string) ([]string, error) {
	var tables []string
	query := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'", schemaName)
	err := pq.db.Raw(query).Scan(&tables).Error
	if err != nil {
		log.Printf("failed to 'ReadAllTables' agr(%s)", schemaName)
		return nil, err
	}
	return tables, err
}

func (pq *PostgreSQL) ExistSchema(schemaName string) (bool, error) {
	var schemaNameResult string
	err := pq.db.Raw("SELECT schema_name FROM information_schema.schemata WHERE schema_name = ?", schemaName).Scan(&schemaNameResult).Error
	if err != nil {
		log.Printf("failed to 'ExistSchema' arg(%s, %v)", schemaName, err)
		return true, err
	} else if schemaNameResult == "" {
		return false, nil
	} else {
		return true, nil
	}
}

func (pq *PostgreSQL) ReadAllSchemas() ([]string, error) {
	var schemas []string
	err := pq.db.Raw("SELECT schema_name FROM information_schema.schemata").Scan(&schemas).Error
	if err != nil {
		log.Printf("failed to 'ReadAllSchemas' (%v)", err)
		return nil, err
	}
	return schemas, nil
}

// users 레코드들을 삭제 후 다시 AddUser
func (pq *PostgreSQL) UpdateUser(users []*auth.User) error {

	var emails strings.Builder
	for i, user := range users {
		buf := fmt.Sprintf("'%s'", user.Email)
		emails.WriteString(buf)
		if len(users)-1 == i {
			continue
		}
		emails.WriteString(",")
	}

	var ids []struct{ ID int }
	if err := pq.db.Raw(fmt.Sprintf(`SELECT ID FROM AUTH.USERS WHERE EMAIL IN (%s)`, emails.String())).
		Scan(&ids).Error; err != nil {
		log.Println("failed to select auth.user.id", err)
		return err
	}

	idArray := make([]int, len(ids))
	for i, id := range ids {
		idArray[i] = id.ID
	}
	if err := pq.db.Delete(&auth.User{}, idArray).Error; err != nil {
		log.Println("failed to delete user", err)
		return err
	}

	if err := pq.AddUser(users); err != nil {
		return err
	}

	return nil
}

func (pq *PostgreSQL) AddUser(users []*auth.User) error {
	for i, user := range users {
		roles := user.Roles
		for j, role := range roles {
			id, err := utils.GetRoleId(role.RoleName)
			if err != nil {
				log.Println(err)
				return err
			}
			users[i].Roles[j].ID = id
		}
	}
	if err := pq.db.Create(users).Error; err != nil {
		log.Println("failed to add user", err)
		return err
	}

	// insert user_roles mapping table
	/*
			명시적으로 맵핑 테이블에 데이터를 추가해야 할 경우, 아래와 같이 쿼리를 작성할 수도 있음.

		var roleNames strings.Builder
		roles := user.Roles
		for i, role := range roles {
			buf := fmt.Sprintf("'%s'", role.RoleName)
			roleNames.WriteString(buf)
			if len(roles)-1 == i {
				continue
			}
			roleNames.WriteString(",")
		}
		if err := pq.db.Exec(fmt.Sprintf(
			`INSERT INTO AUTH.USER_ROLES (USER_ID, ROLE_ID)
				SELECT
					A.ID AS USER_ID,
					B.ID AS ROLE_ID
				FROM
					AUTH.USERS AS A
					CROSS JOIN AUTH.ROLES AS B
				WHERE
					A.USERNAME LIKE '%s' AND B.ROLE_NAME IN (%s)`, user.Username, roleNames.String())).Error; err != nil {
			log.Println("failed to load user_roles", err)
			return err
		}
	*/
	return nil
}
