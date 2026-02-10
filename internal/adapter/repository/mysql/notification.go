package mysql

import (
	"database/sql"
	"log"
	"math"
	"strings"

	"github.com/SevgiF/notification-system/internal/core/notification/domain"
	"github.com/SevgiF/notification-system/pkg/defer_util"
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

func (r *NotificationRepository) GetAllNotification(
	status *int,
	channel, from, to *string,
	limit, offset int,
) (items []domain.Notification, totalCount, totalPage int, err error) {

	var where string
	whereArgs := []any{}

	if status != nil && *status != 0 {
		where += ` AND n.status = ? `
		whereArgs = append(whereArgs, *status)
	}
	if channel != nil && strings.TrimSpace(*channel) != "" {
		ch := strings.TrimSpace(*channel)
		where += ` AND n.channel = ? `
		whereArgs = append(whereArgs, ch)
	}
	if from != nil && strings.TrimSpace(*from) != "" {
		f := strings.TrimSpace(*from)
		where += ` AND n.created_at >= ? `
		whereArgs = append(whereArgs, f)
	}
	if to != nil && strings.TrimSpace(*to) != "" {
		t := strings.TrimSpace(*to)
		where += ` AND n.created_at <= ? `
		whereArgs = append(whereArgs, t)
	}

	listQuery := `
SELECT
	n.id, n.recipient, n.channel, n.content, n.priority,
	s.description,
	n.created_at, n.updated_at, n.deleted_at
FROM notification AS n
INNER JOIN status AS s ON s.code = n.status
WHERE 1=1 ` + where + `
ORDER BY n.created_at DESC
LIMIT ?, ?`

	listArgs := append(append([]any{}, whereArgs...), offset, limit)

	log.Print(listQuery)
	rows, err := r.db.Query(listQuery, listArgs...)
	if err != nil {
		return
	}
	defer defer_util.DeferWithErrorHandling(rows.Close)

	for rows.Next() {
		item := domain.Notification{}
		err = rows.Scan(
			&item.ID, &item.Recipient, &item.Channel, &item.Content, &item.Priority,
			&item.StatusName,
			&item.CreatedAt, &item.UpdatedAt, &item.DeletedAt,
		)
		if err != nil {
			return
		}
		items = append(items, item)
	}

	totalQuery := `SELECT COUNT(*) FROM notification AS n WHERE 1=1 ` + where
	log.Println(totalQuery)
	if err = r.db.QueryRow(totalQuery, whereArgs...).Scan(&totalCount); err != nil {
		return
	}

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

func (r *NotificationRepository) CreateNotification(item domain.Notification) (err error) {
	query, err := r.db.Prepare(`INSERT INTO notification (recipient, channel, content, priority, status) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return
	}
	defer defer_util.DeferWithErrorHandling(query.Close)

	_, err = query.Exec(item.Recipient, item.Channel, item.Content, item.Priority, item.Status)
	if err != nil {
		return
	}
	return
}

func (r *NotificationRepository) GetMetrics() (metrics map[string]int, err error) {
	metrics = make(map[string]int)

	// Queue Depth (Created status = 1, Processing = 5)
	var queueDepth int
	err = r.db.QueryRow(`SELECT COUNT(*) FROM notification WHERE status IN (1, 5)`).Scan(&queueDepth)
	if err != nil {
		return
	}
	metrics["queue_depth"] = queueDepth

	// Success/Failure counts (Sent = 3, Failed = 4)
	rows, err := r.db.Query(`SELECT status, COUNT(*) FROM notification WHERE status IN (3, 4) GROUP BY status`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var status, count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		if status == 3 {
			metrics["success_count"] = count
		} else if status == 4 {
			metrics["failure_count"] = count
		}
	}

	return
}

func (r *NotificationRepository) FetchPendingNotifications(limit int) (items []domain.Notification, err error) {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Select pending notifications with locking
	// Priority order: HIGH > NORMAL > LOW (assuming string comparison works or we map them.
	// Actually 'high' < 'low' alphabetically? No. 'high', 'normal', 'low'.
	// We need custom ordering.
	// FIELD(priority, 'high', 'normal', 'low')

	query := `SELECT id, recipient, channel, content, priority, status, created_at, updated_at, deleted_at 
	          FROM notification 
	          WHERE status = ? 
	          ORDER BY CASE priority 
	              WHEN 'high' THEN 1 
	              WHEN 'normal' THEN 2 
	              WHEN 'low' THEN 3 
	              ELSE 4 END ASC, 
	          created_at ASC 
	          LIMIT ? FOR UPDATE`

	rows, err := tx.Query(query, domain.CREATED, limit)
	if err != nil {
		return
	}
	defer rows.Close()

	ids := []int{}

	for rows.Next() {
		item := domain.Notification{}
		err = rows.Scan(&item.ID, &item.Recipient, &item.Channel, &item.Content, &item.Priority, &item.Status, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt)
		if err != nil {
			return
		}
		items = append(items, item)
		ids = append(ids, item.ID)
	}
	rows.Close()

	if len(ids) > 0 {
		// Update status to 'processing' (we don't have PROCESSING status in domain yet, maybe use a placeholder or just keep them locked?
		// Logic says "Process notifications asynchronously".
		// If we don't update status, next poll will pick them up if transaction commits.
		// We should introduce a PROCESSING status = 5 (for example) or similar.
		// But the status table has: 0 deleted, 1 created, 2 canceled, 3 sent, 4 failed.
		// Let's assume we pick them up and process immediately.
		// But 'FetchPending' implies we want to claim them.
		// If we simply return them, the transaction commits and lock releases, so other workers pick them up.
		// We MUST update them to something intermediate or 'sent'/'failed' after processing.
		// For polling, we typically update to 'PROCESSING' to hide from other workers.
		// Since I can't easily change the domain consts without editing multiple files,
		// I'll add PROCESSING = 5 to domain/entity.go and the database.

		// Let's assume I will add PROCESSING status.

		// However, to be safe with the current 'FetchPending' signature, I should probably return them and let the worker update them?
		// No, race condition.

		// I will update them to a temporary status OR I will hold the row lock?
		// No, I can't hold row lock across the function return because Tx commits.

		// So I MUST update them to a 'PROCESSING' status.
		// I will add 'PROCESSING' to the status table in my plan or just use a magic number for now and update domain later.
		// Let's use 5 for PROCESSING.

		// Logic:
		// 1. Select FOR UPDATE
		// 2. Update status to 5 (PROCESSING)
		// 3. Commit
		// 4. Return items.

		// Construct update query
		// "UPDATE notification SET status = 5 WHERE id IN (...)"

		// We need to handle the list of IDs
		// For simplicity, let's just loop update or use logic.

		updateQuery := `UPDATE notification SET status = 5 WHERE id = ?`
		updateStmt, err := tx.Prepare(updateQuery)
		if err != nil {
			return nil, err
		}
		defer updateStmt.Close()

		for _, id := range ids {
			_, err = updateStmt.Exec(id)
			if err != nil {
				return nil, err
			}
		}
	}

	return
}
