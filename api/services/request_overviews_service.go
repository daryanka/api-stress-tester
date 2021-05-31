package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/daryanka/api-stress-tester/api/domains/domains"
	"github.com/daryanka/api-stress-tester/api/domains/individual_requests"
	"github.com/daryanka/api-stress-tester/api/domains/request_overviews"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/daryanka/api-stress-tester/api/websocket_conn"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type RequestOverviewServiceI interface {
	All(userID int64) ([]request_overviews.RequestOverview, utils.RestErrI)
	Delete(id, userID int64) utils.RestErrI
	GetSingle(id, userID int64) (*request_overviews.RequestOverview, utils.RestErrI)
	Create(r request_overviews.NewRequest) (int64, utils.RestErrI)
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

func (i *requestOverviewService) Create(r request_overviews.NewRequest) (int64, utils.RestErrI) {
	// Max num of requests per minute is 1000
	if (r.NumRequests / r.Time) > 1000 {
		return 0, utils.NewBadRequest("The maximum allowed requests per minute is 1000")
	}

	// Check domain belongs to user
	domainInfo, err := domains.DomainsDao.GetSingle(r.DomainID, r.UserID)
	if err != nil {
		return 0, utils.StandardInternalServerError()
	}

	overview := request_overviews.RequestOverview{
		UserID:              r.UserID,
		DomainID:            r.DomainID,
		Endpoint:            r.Endpoint,
		Method:              r.Method,
		Payload:             r.Payload,
		Time:                r.Time,
		NumRequests:         r.NumRequests,
		SuccessfulReq:       0,
		FailedReq:           0,
		AverageResponseTime: 0,
		Status:              request_overviews.StatusInProgress,
	}

	// Create request overview
	id, err := request_overviews.RequestOverviewDao.Create(overview)

	if err != nil {
		return 0, utils.StandardInternalServerError()
	}

	overview.ID = id

	go startNewRequests(overview, domainInfo.DomainURL)

	return id, nil
}

func startNewRequests(overview request_overviews.RequestOverview, baseURL string) {
	// TODO calculate time needed to wait between requests
	waitTime := time.Second * 1

	resultsChan := make(chan individual_requests.IndividualRequest, 100)
	resultsSlice := []individual_requests.IndividualRequest{}
	done := make(chan struct{})
	var wg sync.WaitGroup
	websocket_conn.AddDoneChan(overview.ID, done)

	go func() {
		for el := range resultsChan {
			el.RequestOverviewId = overview.ID
			resultsSlice = append(resultsSlice, el)
			wg.Done()
			go func(r request_overviews.RequestOverview, ir individual_requests.IndividualRequest) {
				b, _ := json.Marshal(websocket_conn.MessageType{
					Type:   websocket_conn.INDIVIDUAL_REQUEST,
					Status: 200,
					Info: websocket_conn.CompletedIndividualRequest{
						RequestID:  r.ID,
						TimeTaken:  ir.TimeTaken,
						StatusCode: ir.StatusCode,
					},
				})
				websocket_conn.WriteMessage(r.UserID, b)
			}(overview, el)
		}
	}()

	fullURL := utils.MergeBaseURLAndEndpoint(baseURL, overview.Endpoint)

	cancelled := false
	for i := 0; i < overview.NumRequests; i++ {
		time.Sleep(waitTime)
		select {
		case <-done:
			cancelled = true
			break
		default:
			wg.Add(1)
			go makeRequest(fullURL, overview.Method, overview.Payload, resultsChan)
		}
		if cancelled {
			break
		}
	}
	wg.Wait()

	if !cancelled {
		websocket_conn.RemoveDoneChan(overview.ID)
		close(done)
	}
	close(resultsChan)

	// Generate Averages and save to database
	err := individual_requests.IndividualRequestsDao.Create(resultsSlice)
	if err != nil {
		b, _ := json.Marshal(websocket_conn.MessageType{
			Type:   websocket_conn.REQUEST_FAILED,
			Status: 500,
			Info: websocket_conn.RequestUpdate{
				RequestID: overview.ID,
				Message:   "Server error",
			},
		})
		websocket_conn.WriteMessage(overview.UserID, b)
		return
	}

	// Successful requests
	successful := 0
	failed := 0
	var averageResponseTime float64 = 0
	for _, el := range resultsSlice {
		if el.StatusCode >= 200 && el.StatusCode < 300 {
			successful++
		} else {
			failed++
		}
		averageResponseTime += float64(el.TimeTaken)
	}
	averageResponseTime = averageResponseTime / float64(len(resultsSlice))

	overview.SuccessfulReq = successful
	overview.FailedReq = failed
	overview.AverageResponseTime = averageResponseTime
	if cancelled {
		overview.Status = request_overviews.StatusCancelled
	} else {
		overview.Status = request_overviews.StatusComplete
	}

	err = request_overviews.RequestOverviewDao.UpdateRequestResults(overview)
	if err != nil {
		b, _ := json.Marshal(websocket_conn.MessageType{
			Type:   websocket_conn.REQUEST_FAILED,
			Status: 500,
			Info: websocket_conn.RequestUpdate{
				RequestID: overview.ID,
				Message:   "Server error",
			},
		})
		websocket_conn.WriteMessage(overview.UserID, b)
		return
	}

	b, _ := json.Marshal(websocket_conn.MessageType{
		Type:   websocket_conn.REQUEST_COMPLETE,
		Status: 200,
		Info: websocket_conn.RequestUpdate{
			RequestID: overview.ID,
			Message:   "All requests have been made you can now check the results",
		},
	})
	websocket_conn.WriteMessage(overview.UserID, b)
}

func makeRequest(fullURL, method string, jsonBody *string, resChan chan<- individual_requests.IndividualRequest) {
	// TODO Test timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	var body io.Reader

	if jsonBody != nil {
		body = strings.NewReader(*jsonBody)
	}

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return
	}

	if jsonBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	before := time.Now()
	res, err := client.Do(req)
	timeTaken := time.Since(before).Milliseconds()
	if timeTaken == 0 {
		timeTaken = 1
	}
	if err != nil {
		// return status code 0 meaning time out and time to 10 seconds
		fmt.Println("error", err.Error())
		resChan <- individual_requests.IndividualRequest{
			StatusCode: 0,
			TimeTaken:  10000,
		}
		return
	}
	_ = res.Body.Close()

	// Add status code and time
	resChan <- individual_requests.IndividualRequest{
		StatusCode: res.StatusCode,
		TimeTaken:  timeTaken,
	}
}
