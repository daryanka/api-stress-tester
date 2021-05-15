package domains

import (
	"database/sql"
	"github.com/daryanka/api-stress-tester/api/clients"
	"github.com/daryanka/api-stress-tester/api/utils"
)

type DomainsDaoI interface {

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

func (d *domainsDao) Create(domain *CreateDomain) error {
	_, err := clients.DB.Exec(queryAddDomain, domain.DomainURL, domain.Token, domain.Verified, domain.UserID)
	if err != nil {
		utils.Logger.Error("error creating domain ", err)
	}
	return err
}

func (d *domainsDao) Confirm(id int64) error {
	_, err := clients.DB.Exec(queryVerifyDomain, id)
	if err != nil {
		utils.Logger.Error("error verifying domain ", err)
	}
	return err
}

func (d *domainsDao) Delete(id, user_id int64) error {
	_, err := clients.DB.Exec(queryDelete, id, user_id)
	if err != nil {
		utils.Logger.Error("error deleting domain ", err)
	}

	return err
}