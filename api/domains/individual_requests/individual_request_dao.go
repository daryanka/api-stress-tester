package individual_requests

import (
	"database/sql"
	"github.com/daryanka/api-stress-tester/api/clients"
	"github.com/daryanka/api-stress-tester/api/utils"
)

type IndividualRequestDaoI interface {
	GetAll(reqID int64) (res []IndividualRequest, err error)
	Create(data IndividualRequest) (err error)
}

type individualRequestsDao struct {}

var IndividualRequestsDao IndividualRequestDaoI = &individualRequestsDao{}

func (i *individualRequestsDao) GetAll(reqID int64) (res []IndividualRequest, err error) {
	err = clients.DB.Select(&res, queryGetAll, reqID)
	if err != nil && err != sql.ErrNoRows {
		utils.Logger.Error("error getting individual requests ", err)
	}
	return
}

func (i *individualRequestsDao) Create(data IndividualRequest) (err error) {
	_, err = clients.DB.Exec(queryCreate, data.RequestOverviewId, data.StatusCode, data.TimeTaken)
	if err != nil {
		utils.Logger.Error("error creating individual request ", err)
	}
	return
}