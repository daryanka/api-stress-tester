package services

import (
	"database/sql"
	"github.com/daryanka/api-stress-tester/api/domains/domains"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/google/uuid"
	"net/url"
)

type DomainServiceI interface {
	All(userID int64) ([]domains.Domain, utils.RestErrI)
	Create(req domains.CreateDomain) (int64, utils.RestErrI)
	Verify(id int64) utils.RestErrI
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
	req.DomainURL = validURL

	// Create token
	req.Token = uuid.New().String()

	id, err := domains.DomainsDao.Create(&req)
	if err != nil {
		return 0, utils.StandardInternalServerError()
	}
	return id, nil
}

func (d *domainService) Verify(id int64) utils.RestErrI {
	// TODO make API request, and check token is valid
	err := domains.DomainsDao.VerifyToken(id)
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
