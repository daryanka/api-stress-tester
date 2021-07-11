package domains

type Domain struct {
	ID          int64  `json:"id" db:"id"`
	DomainURL   string `json:"domain_url" db:"domain_url"`
	UserID      int64  `json:"user_id" db:"user_id"`
	Verified    int    `json:"verified" db:"verified"`
	Token       string `json:"token" db:"token"`
	NumRequests int    `json:"num_requests" db:"num_requests"`
}

type CreateDomain struct {
	DomainURL string `json:"domain_url" validate:"required,max=254"`
	UserID    int64  `json:"user_id"`
	Token     string `json:"token"`
	Verified  int    `json:"verified"`
}

type TokenVerify struct {
	Endpoint string `json:"endpoint" validate:"required"`
	ID       int64  `json:"id" validate:"required"`
}
