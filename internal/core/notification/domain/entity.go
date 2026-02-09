package domain

type Notification struct {
	ID        int
	Recipient string
	Channel   string
	Content   string
	Priority  string
	Status    int
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

type Status struct {
	Code        int
	Description string
}

const (
	HIGH   = "high"
	NORMAL = "normal"
	LOW    = "low"

	//statuses
	DELETED  = 0
	CREATED  = 1
	CANCELED = 2

	//rules
	MAX_LIMIT     = 100
	MIN_LIMIT     = 10
	DEFAULT_LIMIT = 20
)
