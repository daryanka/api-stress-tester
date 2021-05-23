package services

import (
	"database/sql"
	"github.com/daryanka/api-stress-tester/api/domains/individual_requests"
	"github.com/daryanka/api-stress-tester/api/domains/request_overviews"
	"github.com/daryanka/api-stress-tester/api/utils"
)

type RequestOverviewServiceI interface {
	All(userID int64) ([]request_overviews.RequestOverview, utils.RestErrI)
	Delete(id, userID int64) utils.RestErrI
	GetSingle(id, userID int64) (*request_overviews.RequestOverview, utils.RestErrI)
}

type requestOverviewService struct{}

var RequestOverviewService RequestOverviewServiceI = &requestOverviewService{}

func (i *requestOverviewService) All(userID int64) ([]request_overviews.RequestOverview, utils.RestErrI) {
	res, err := request_overviews.RequestOverviewDao.GetAll(userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, utils.StandardInternalServerError()
	}
	if res == nil {
		return []request_overviews.RequestOverview{}, nil
	}
	return res, nil
}

func (i *requestOverviewService) Delete(id, userID int64) utils.RestErrI {
	err := request_overviews.RequestOverviewDao.Delete(userID, id)
	if err != nil {
		return utils.StandardInternalServerError()
	}
	return nil
}

func (i *requestOverviewService) GetSingle(id, userID int64) (*request_overviews.RequestOverview, utils.RestErrI) {
	res, err := request_overviews.RequestOverviewDao.GetSingle(userID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewBadRequest("Request not found")
		}
		return nil, utils.StandardInternalServerError()
	}

	// Get requests individual requests
	single, err := individual_requests.IndividualRequestsDao.GetAll(res.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, utils.StandardInternalServerError()
	}
	if len(single) == 0 {
		res.IndividualRequests = []individual_requests.IndividualRequest{}
	} else {
		res.IndividualRequests = single
	}

	return &res, nil
}