package repository

import (
	"autocall/common/log"
	"autocall/common/model"
	IMySql "autocall/internal/sqldb/mysql"
)

type (
	IUser interface {
	}
	User struct{}
)

var UserRepo IUser

func NewUserRepo() User {
	repo := User{}
	return repo
}

func (repo *User) CreateUser(entry model.User) (interface{}, error) {
	err := IMySql.MySqlConnector.GetConn().Model(&model.User{}).Create(&entry).Error
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (repo *User) GetUserByApiKey(apiKey string) (interface{}, error) {
	log.Info("UserRepository", "GetUserByApiKey", apiKey)
	var user model.User
	err := IMySql.MySqlConnector.GetConn().Model(&model.User{}).Where("api_key = ?", apiKey).First(&user).Error
	if err != nil {
		log.Error("UserRepository", "GetUserByApiKey - error", err)
		return nil, err
	}
	return user, nil
}

/**
* Get info user, include level, group_id
 */
func (repo *User) SelectUserInGroup(groupID string, limit, offset int) ([]map[string]interface{}, int64, error) {
	log.Info("UserRepository", "GetUserInGroup", groupID)
	var users []map[string]interface{}
	query := IMySql.MySqlConnector.GetConn().Debug().Model(&model.User{}).
		Joins("INNER JOIN group_user ON group_user.user_id=user.id").
		Where("group_user.group_id = ?", groupID).
		Limit(limit).
		Offset(offset).
		Group("user.id").
		Select("user.id, username, api_key, level").
		Find(&users)
	total := query.RowsAffected
	if query.Error != nil {
		log.Error("UserRepository", "GetUserInGroup - error", query.Error)
		return nil, 0, query.Error
	}
	// data := []map[string]interface{}{}
	// for _, value := range users {
	// 	data = append(data, value)
	// }
	return users, total, nil
}
