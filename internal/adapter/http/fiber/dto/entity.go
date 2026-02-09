package dto

type NotificationRequest struct {
	Recipient string `json:"recipient"`
	Channel   string `json:"channel"`
	Content   string `json:"content"`
	Priority  string `json:"priority"`
}

type NotificationResponse struct {
	ID         int    `json:"id"`
	Recipient  string `json:"recipient"`
	Channel    string `json:"channel"`
	Content    string `json:"content"`
	Priority   string `json:"priority"`
	Status     int    `json:"status"`
	StatusName string `json:"statusName"`
}

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type Response struct {
	Message    string      `json:"message"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Filter     *Filter     `json:"filter,omitempty"`
	Data       any         `json:"data,omitempty"`
}

type Pagination struct {
	Limit       int `json:"limit"`
	CurrentPage int `json:"currentPage"`
	TotalCount  int `json:"totalCount"`
	TotalPage   int `json:"totalPage"`
}

type Filter struct {
	From    *string `json:"from"`
	To      *string `json:"to"`
	Status  *int    `json:"status"`
	Channel *string `json:"channel"`
}
