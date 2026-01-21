package service

import (
	"log"
	"os"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
)

const batchSize = 5

func RunReminderCron() {
	today := time.Now().Format("2006-01-02")
	log.Println("reminder cron: starting reminder")
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
			emailData := entity.SubscriptionEmailData{
				ID:              "3d08f42e-c76f-4d77-b096-31acbe7506c7",
				OrderID:         1,
				UserName:        "Raihan",
				ProductName1:    "Iron Nexus Plastic",
				ProductName2:    "A",
				Subscription:    "2 Bulan",
				Status:          "✅ Aktif",
				OrderDate:       "2026-01-15 09:18:42",
				PaymentDeadline: "2026-01-19",
				TargetEmail:     "mraihanzhafran.14@gmail.com",
				Domain:          os.Getenv("WEBSITE_URL"),
			}

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
			if checkSchedule(r.NextRemove, today) {
				log.Println(
					"REMOVE: ",
					r.ID,
					r.UserID,
					r.ProductID,
					r.ProductVariantID,
				)

				_ = model.UpdateLastSentAt(r.ID, time.Now())
				_ = model.UpdateIsOver(r.ID)
			} else if checkSchedule(r.NextWarningSend, today) {
				log.Println(
					"WARNING: ",
					r.ID,
					r.UserID,
					r.ProductID,
					r.ProductVariantID,
				)

				_ = model.UpdateLastSentAt(r.ID, time.Now())
			} else if checkSchedule(r.NextSend, today) {
				log.Println(
					"PROCESS: ",
					r.ID,
					r.UserID,
					r.ProductID,
					r.ProductVariantID,
				)

				_ = model.UpdateLastSentAt(r.ID, time.Now())
			}

			htmlBody, err := RenderSubscriptionEmail(emailData, "E:/GIU/Devel/go_app/subcommerce-backend/internal/templates/email_schedule_reminder.html")
			if err != nil {
				log.Fatal(err)
			}

			sender := NewBrevoSender()
			_ = sender.SendMail(
				"mraihanzhafran.14@gmail.com",
				"Pengingat Langganan Product Subcommerce",
				htmlBody,
			)
		}
	}
}

func checkSchedule(t time.Time, today string) bool {
	return t.Format("2006-01-02") == today
}
