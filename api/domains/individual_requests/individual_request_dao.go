package individual_requests

import (
	"database/sql"
	"github.com/daryanka/api-stress-tester/api/clients"
	"github.com/daryanka/api-stress-tester/api/utils"
	"strings"
)

type IndividualRequestDaoI interface {
	GetAll(reqID int64) (res []IndividualRequest, err error)
	Create(data []IndividualRequest) (err error)
}

type individualRequestsDao struct{}

var IndividualRequestsDao IndividualRequestDaoI = &individualRequestsDao{}

func (i *individualRequestsDao) GetAll(reqID int64) (res []IndividualRequest, err error) {
	err = clients.DB.Select(&res, queryGetAll, reqID)
	if err != nil && err != sql.ErrNoRows {
		utils.Logger.Error("error getting individual requests ", err)
	}
	return
}

func (i *individualRequestsDao) Create(data []IndividualRequest) (err error) {
	// Convert individual request into multi requests each with 1000 rows
	args := [][]interface{}{}

	current := []interface{}{}
	for index, el := range data {
		if (index+1)%1000 == 0 {
			current = append(current, el.RequestOverviewId, el.StatusCode, el.TimeTaken)
			args = append(args, current)
			current = []interface{}{}
		} else {
			current = append(current, el.RequestOverviewId, el.StatusCode, el.TimeTaken)
		}

		if index == len(data)-1 && len(current) > 0 {
			args = append(args, current)
		}
	}

	for _, el := range args {
		var str []string
		for num := 0; num < len(el) / 3; num++ {
			str = append(str, "(?, ?, ?)")
		}
		_, err = clients.DB.Exec(queryCreate+strings.Join(str, ",")+";", el...)
		if err != nil {
			utils.Logger.Error("error creating individual request ", err)
		}
	}

	return
}
