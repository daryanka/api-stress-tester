package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/daryanka/api-stress-tester/api/domains/domains"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
)

type DomainServiceI interface {
	All(userID int64) ([]domains.Domain, utils.RestErrI)
	Create(req domains.CreateDomain) (int64, utils.RestErrI)
	Verify(endpoint string, id, userID int64) utils.RestErrI
	Delete(id int64, userID int64) utils.RestErrI
}

type domainService struct{}

var DomainService DomainServiceI = &domainService{}

func (d *domainService) All(userID int64) ([]domains.Domain, utils.RestErrI) {
	res, err := domains.DomainsDao.GetAll(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []domains.Domain{}, nil
		}
		return nil, utils.StandardInternalServerError()
	}
	if res == nil {
		return []domains.Domain{}, nil
	}
	return res, nil
}

func (d *domainService) Create(req domains.CreateDomain) (int64, utils.RestErrI) {
	// Verify if URL is valid
	parsedURL, err := url.ParseRequestURI(req.DomainURL)
	if err != nil {
		return 0, utils.NewUnprocessableEntity("Invalid URL provided")
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return 0, utils.NewUnprocessableEntity("Invalid URL provided")
	}

	validURL := parsedURL.Scheme + "://" + parsedURL.Hostname()
	if parsedURL.Port() != "" {
		validURL = validURL + ":" + parsedURL.Port()
	}
	req.DomainURL = validURL

	// Create token
	req.Token = uuid.New().String()

	id, err := domains.DomainsDao.Create(&req)
	if err != nil {
		return 0, utils.StandardInternalServerError()
	}
	return id, nil
}

func (d *domainService) Verify(endpoint string, id, userId int64) utils.RestErrI {
	fmt.Println("here", id, userId)
	domain, err := domains.DomainsDao.GetSingle(id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.StandardUnauthorized()
		}
		return utils.StandardInternalServerError()
	}

	fmt.Println(domain)

	if e := makeVerifyDomainRequest(domain.DomainURL + endpoint, domain.Token); e != nil {
		return e
	}

	err = domains.DomainsDao.VerifyToken(id)
	if err != nil {
		return utils.StandardInternalServerError()
	}
	return nil
}

func (d *domainService) Delete(id int64, userID int64) utils.RestErrI {
	err := domains.DomainsDao.Delete(id, userID)
	if err != nil {
		return utils.StandardInternalServerError()
	}
	return nil
}

func makeVerifyDomainRequest(url, token string) utils.RestErrI {
	type Data struct {
		Token string `json:"token"`
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return utils.NewBadRequest("Unable to make request", "INVALID URL")
	}

	res, err := client.Do(req)

	if err != nil {
		return utils.StandardInternalServerError("UNABLE TO MAKE REQUEST")
	}

	b, err := ioutil.ReadAll(res.Body)

	body := Data{}
	if err = json.Unmarshal(b, &body); err != nil {
		return utils.NewBadRequest("Unable to read response body please check that the token is being returned in JSON format")
	}

	if body.Token != token {
		return utils.NewBadRequest("Invalid token found in response")
	}
	return nil
}
