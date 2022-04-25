package repository

import (
	"autocall/common/log"
	"autocall/common/model"
	IMySql "autocall/internal/sqldb/mysql"
)

type (
	IContact interface {
		SelectContactByID(id string) (interface{}, error)
	}
	Contact struct{}
)

var ContactRepo IContact

func NewContactRepo() Contact {
	repo := Contact{}
	return repo
}

// func (repo *Contact) SelectContactByID(id string) (interface{}, error) {
// 	contact := new(model.Contact)
// 	err := IMySql.MySqlConnector.GetConn().Debug().Model(&model.Contact{}).Where("id = ?", id).First(contact).Error
// 	if err != nil {
// 		log.Error("ContactRepository", "SelectContactByID - error", err)
// 		return nil, err
// 	}
// 	log.Info("ContactRepository", "SelectContactByID - success", contact)
// 	return contact, nil
// }

func (repo *Contact) SelectContactByID(id, userID, groupID, level string) (interface{}, error) {
	// userID, groupID, level := mdw
	log.Info("ContactRepository", "SelectContactByID - start", "")
	contact := new(model.Contact)
	row := map[string]interface{}{}
	query := IMySql.MySqlConnector.GetConn().Debug().Model(&model.Contact{}).
		Joins("INNER JOIN user ON user.id=contact_call_form.user_id").
		Joins("INNER JOIN group_user ON group_user.user_id=user.id")
	if level <= "2" && level != "0" {
		query = query.Where("user.level <= ?", level).
			Where("group_user.group_id=?", groupID)
	} else if level == "0" {
		query = query.Where("user.level = ?", level).
			Where("user.id = ?", userID)
	}
	query = query.Select("contact_call_form.*")
	err := query.Find(&row).Error
	if err != nil {
		log.Error("ContactRepository", "SelectContactByID - error", err)
		return nil, err
	}
	log.Info("ContactRepository", "SelectContactByID - success", contact)
	return row, nil
}
