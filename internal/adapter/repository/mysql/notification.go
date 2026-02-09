package mysql

import (
	"database/sql"
	"github.com/SevgiF/notification-system/internal/core/notification/domain"
	"github.com/SevgiF/notification-system/pkg/defer_util"
	"math"
)

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

type NotificationRepository struct {
	db *sql.DB
}

func (r *NotificationRepository) GetStatusList() (statuses []domain.Status, err error) {
	rows, err := r.db.Query(`SELECT code, description FROM status ORDER BY code ASC`)
	if err != nil {
		return
	}
	defer defer_util.DeferWithErrorHandling(rows.Close)

	for rows.Next() {
		var status domain.Status
		err = rows.Scan(&status.Code, &status.Description)
		if err != nil {
			return
		}
		statuses = append(statuses, status)
	}
	return
}

func (r *NotificationRepository) GetAllNotification(status *int, channel, from, to *string, limit, offset int) (items []domain.Notification, totalCount, totalPage int, err error) {
	var conditions string
	args := []any{}
	if status != nil {
		conditions += ` AND n.status = ? `
		args = append(args, *status)
	}
	if channel != nil {
		conditions += ` AND n.channel = ? `
		args = append(args, *channel)
	}
	if from != nil {
		conditions += ` AND n.created_at >= ? `
		args = append(args, *from)
	}
	if to != nil {
		conditions += ` AND n.created_at <= ? `
		args = append(args, *to)
	}
	conditions += ` LIMIT ?, ?`
	args = append(args, offset, limit)

	query := `SELECT n.id, n.recipient, n.channel, n.content, n.priority, s.description, n.created_at, n.updated_at, n.deleted_at
							FROM notification AS n
							JOIN status AS s ON n.status=s.code
							WHERE 1=1 ` + conditions
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return
	}
	defer defer_util.DeferWithErrorHandling(rows.Close)

	for rows.Next() {
		item := domain.Notification{}
		err = rows.Scan(&item.ID, &item.Recipient, &item.Channel, &item.Content, &item.Priority, &item.Status, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt)
		if err != nil {
			return
		}
		items = append(items, item)
	}

	totalQuery := `SELECT COUNT(*) FROM notification WHERE 1=1 ` + conditions
	err = r.db.QueryRow(totalQuery, args...).Scan(&totalCount)

	totalPageCalc := float64(totalCount) / float64(limit)
	totalPage = int(math.Ceil(totalPageCalc))
	if totalPage > 0 {
		totalPage--
	}
	return
}

func (r *NotificationRepository) UpdateNotification(id int, status int) (err error) {
	query, err := r.db.Prepare(`UPDATE notification SET status=? WHERE id=?`)
	if err != nil {
		return
	}
	defer defer_util.DeferWithErrorHandling(query.Close)
	_, err = query.Exec(status, id)
	if err != nil {
		return
	}
	return
}

func (r *NotificationRepository) GetNotification(id int) (item domain.Notification, err error) {
	err = r.db.QueryRow(`SELECT id, recipient, channel, content, priority, status, created_at, updated_at, deleted_at FROM notification WHERE id=?`, id).
		Scan(&item.ID, &item.Recipient, &item.Channel, &item.Content, &item.Priority, &item.Status, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt)
	if err != nil {
		return
	}
	return
}

func (r *NotificationRepository) CreateNotification(tx *sql.Tx, item domain.Notification) (err error) {
	query, err := tx.Prepare(`INSERT INTO notification (reciepent, channel, content, priority) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return
	}
	defer defer_util.DeferWithErrorHandling(query.Close)

	_, err = query.Exec(item.Recipient, item.Channel, item.Content, item.Priority)
	if err != nil {
		return
	}
	return
}

func (r *NotificationRepository) WithTransaction() (*sql.Tx, error) {
	return r.db.Begin()
}
