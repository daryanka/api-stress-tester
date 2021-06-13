package request_overviews

import (
	"database/sql"
	"github.com/daryanka/api-stress-tester/api/clients"
	"github.com/daryanka/api-stress-tester/api/utils"
)

type RequestOverviewDaoI interface {
	GetAll(userID int64) (res []RequestOverview, err error)
	GetSingle(userID, id int64) (res RequestOverview, err error)
	Delete(userID, id int64) (err error)
	UpdateRequestResults(data RequestOverview) (err error)
	Create(data RequestOverview) (id int64, err error)
	UpdateStatus(status int, id int64) (err error)
}

type requestOverviewDao struct{}

var RequestOverviewDao RequestOverviewDaoI = &requestOverviewDao{}

func (i *requestOverviewDao) GetAll(userID int64) (res []RequestOverview, err error) {
	err = clients.DB.Select(&res, queryGetAllRequests, userID)
	if err != nil && err != sql.ErrNoRows {
		utils.Logger.Error("error getting all requests ", err)
	}
	return
}

func (i *requestOverviewDao) GetSingle(userID, id int64) (res RequestOverview, err error) {
	err = clients.DB.Get(&res, queryGetSingle, userID, id)
	if err != nil && err != sql.ErrNoRows {
		utils.Logger.Error("error getting single requests ", err)
	}
	return
}

func (i *requestOverviewDao) Delete(userID, id int64) (err error) {
	_, err = clients.DB.Exec(queryDelete, userID, id)
	if err != nil {
		utils.Logger.Error("error deleting request ", err)
	}
	return
}

func (i *requestOverviewDao) Create(data RequestOverview) (id int64, err error) {
	res, err := clients.DB.Exec(queryCreate,
		data.UserID,
		data.ReqName,
		data.DomainID,
		data.Endpoint,
		data.Method,
		data.Payload,
		data.Time,
		data.NumRequests,
		data.SuccessfulReq,
		data.FailedReq,
		data.AverageResponseTime,
		data.Status,
	)
	if err != nil {
		utils.Logger.Error("error deleting request ", err)
		return
	}
	id, err = res.LastInsertId()
	if err != nil {
		utils.Logger.Error("error getting request id ", err)
	}
	return
}

func (i *requestOverviewDao) UpdateRequestResults(data RequestOverview) (err error) {
	_, err = clients.DB.Exec(queryUpdateResults,
		data.SuccessfulReq,
		data.FailedReq,
		data.AverageResponseTime,
		data.Status,
		data.ID,
	)
	if err != nil {
		utils.Logger.Error("error updating request results ", err)
	}
	return
}

func (i *requestOverviewDao) UpdateStatus(status int, id int64) (err error) {
	_, err = clients.DB.Exec(queryUpdateStatus, status, id)
	if err != nil {
		utils.Logger.Error("error updating request status ", err)
	}
	return
}
