package domains

type Domain struct {
	ID        int64  `json:"id" db:"id"`
	DomainURL string `json:"domain_url" db:"domain_url"`
	UserID    int64  `json:"user_id" db:"user_id"`
	Verified  string `json:"verified" db:"verified"`
	Token     string `json:"token" db:"token"`
}

type CreateDomain struct {
	DomainURL string `json:"domain_url" validate:"required"`
	UserID    int64  `json:"user_id"`
	Token     string `json:"token"`
	Verified   string `json:"verified"`
}
