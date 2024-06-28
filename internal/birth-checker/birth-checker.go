package birthchecker

import (
	"fmt"
	"log"
	"time"

	"github.com/fishmanDK/rutube-test-task/config"
	"github.com/jmoiron/sqlx"
)

const (
	subjectDeleteBanner = "deleteRequest"
)

type event struct {
	id        int64 `db:"id"`
	bannerID  int64 `db:"bannerID"`
	tagID     int64 `db:"tagID"`
	featureID int64 `db:"featureID"`
}

type Checker struct {
	db       *sqlx.DB
}

func NewChecker(cfg config.ConfigDB) (*Checker, error) {
	const op = "event_checker.NewChecker"
	db, err := sqlx.Open("postgres", cfg.ToString())
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &Checker{
		db:       db,
	}, nil
}

func (c *Checker) Run(timeTicker time.Duration) error {
	const op = "event_checker.CheckExpiredRequests"

	ticker := time.NewTicker(timeTicker)
	defer ticker.Stop()

	selectQuery := `
			SELECT id, user_name, email 
			FROM users 
			WHERE birth_date == $1;
	`

	deleteQuery := `
			DELETE FROM deletion_requests
			WHERE id = $1;
	`

	for range ticker.C {

		rows, err := c.db.Query(selectQuery)
		if err != nil {
			log.Println("Error querying database:", err)
			continue
		}
		defer rows.Close()

		for rows.Next() {
			var e event
			if err := rows.Scan(&e.id, &e.bannerID, &e.tagID, &e.featureID); err != nil {
				log.Println("Error scanning row:", err)
				continue
			}
			fmt.Println(e)

			_, err = c.db.Exec(deleteQuery, e.id)
			if err != nil {
				log.Printf("failed delete reqest deletion_requests.id = %d: %v", e.id, err)
			}
		}
	}

	return nil
}