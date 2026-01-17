package service

import (
	"log"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/model"
)

const batchSize = 2

func RunReminderCron() {
	today := time.Now().Format("2006-01-02")
	for {
		rows, err := model.FetchReminderBatch(today, batchSize)
		if err != nil {
			log.Println("fetch error:", err)
			return
		}

		if len(rows) == 0 {
			log.Println("reminder cron: no more data")
			return
		}

		for _, r := range rows {
			// 1️⃣ If is_over → delete
			if r.IsOver {
				_ = model.DeleteReminder(r.ID)
				continue
			}

			// 2️⃣ If already processed today → skip
			if r.LastSentAt.Format("2006-01-02") == today {
				continue
			}

			// 3️⃣ Check schedule dates
			if checkSchedule(r.NextSend, today) {
				log.Println(
					"PROCESS:",
					r.ID,
					r.UserID,
					r.ProductID,
					r.ProductVariantID,
				)

				_ = model.UpdateLastSentAt(r.ID, time.Now())
			} else if checkSchedule(r.NextWarningSend, today) {
				log.Println(
					"WARNING:",
					r.ID,
					r.UserID,
					r.ProductID,
					r.ProductVariantID,
				)

				_ = model.UpdateLastSentAt(r.ID, time.Now())
			} else if checkSchedule(r.NextRemove, today) {
				log.Println(
					"REMOVE:",
					r.ID,
					r.UserID,
					r.ProductID,
					r.ProductVariantID,
				)

				_ = model.UpdateLastSentAt(r.ID, time.Now())
			}
		}
	}
}

func checkSchedule(t time.Time, today string) bool {
	return t.Format("2006-01-02") == today
}
