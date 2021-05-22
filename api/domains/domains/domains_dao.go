package domains

import (
	"database/sql"
	"fmt"
	"github.com/daryanka/api-stress-tester/api/clients"
	"github.com/daryanka/api-stress-tester/api/utils"
)

type DomainsDaoI interface {
	GetAll(userID int64) (res []Domain, err error)
	Create(domain *CreateDomain) (int64, error)
	VerifyToken(id int64) error
	Delete(id, userId int64) error
	GetSingle(id int64, userID int64) (res Domain, err error)
}

type domainsDao struct{}

var DomainsDao DomainsDaoI = &domainsDao{}

func (d *domainsDao) GetAll(userID int64) (res []Domain, err error) {
	err = clients.DB.Select(&res, queryGetAll, userID)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.Logger.Error("error getting all domains ", err)
		}
	}
	return
}

func (d *domainsDao) GetSingle(id int64, userID int64) (res Domain, err error) {
	err = clients.DB.Get(&res, queryGetSingle, userID, id)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.Logger.Error("error getting all domains ", err)
		}
	}
	return
}

func (d *domainsDao) Create(domain *CreateDomain) (int64, error) {
	res, err := clients.DB.Exec(queryAddDomain, domain.DomainURL, domain.Token, domain.Verified, domain.UserID)
	if err != nil {
		utils.Logger.Error("error creating domain ", err)
		return 0, nil
	}
	id, err := res.LastInsertId()
	if err != nil {
		utils.Logger.Error("error getting created domain id ", err)
	}
	return id, err
}

func (d *domainsDao) VerifyToken(id int64) error {
	fmt.Println("running", id)
	_, err := clients.DB.Exec(queryVerifyDomain, id)
	if err != nil {
		utils.Logger.Error("error verifying domain ", err)
	}
	return err
}

func (d *domainsDao) Delete(id, userId int64) error {
	_, err := clients.DB.Exec(queryDelete, id, userId)
	if err != nil {
		utils.Logger.Error("error deleting domain ", err)
	}

	return err
}