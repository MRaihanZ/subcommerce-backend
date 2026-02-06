package service

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func errorUpdateReminderSchedules() {
	log.Println("❌ Error update data in reminder_schedules table function")
	message := recover()
	log.Println("❌ ERROR: ", message)
}

func errorSubscriptionReminderEmail() {
	log.Println("❌ Error parsing html file function")
	message := recover()
	log.Println("ERROR: ", message)
}

const batchSize = 5

func RunReminderCron() {
	today := time.Now().Format("2006-01-02")
	var caser cases.Caser
	var err error
	var rows []entity.SubscriptionData
	for {
		rows, err = model.FetchReminderBatch(today, batchSize)
		if err != nil {
			log.Println("fetch error:", err)
			return
		}

		if len(rows) == 0 {
			log.Println("Reminder Cron: no more data")
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

			caser = cases.Title(language.MustParse("id-ID"))

			emailData := entity.SubscriptionEmailData{
				OrderID:         r.OrderPrettyID,
				SubscriptionID:  r.ID,
				UserName:        r.UName,
				ProductName1:    r.PName,
				ProductName2:    r.PVName,
				Subscription:    strconv.Itoa(r.PVInterval) + " " + caser.String(r.IName),
				OrderDate:       r.CreatedAt.Format("02-01-2006 15:04:05"),
				PaymentDeadline: r.NextRemove.Format("02-01-2006"),
				TargetEmail:     r.UEmail,
				Domain:          os.Getenv("WEBSITE_URL"),
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
				emailData.Status = "🛑 Jatuh Tempo"

				err = model.UpdateIsOver(r.ID)
				if err != nil {
					defer errorUpdateReminderSchedules()
					errMessage := "Failed to update data in table reminder_schedules: " + err.Error()
					panic(errMessage)
				}
			} else if checkSchedule(r.NextWarningSend, today) {
				log.Println(
					"WARNING: ",
					r.ID,
					r.UserID,
					r.ProductID,
					r.ProductVariantID,
				)
				emailData.Status = "⚠️ Akan Berakhir"
			} else if checkSchedule(r.NextSend, today) {
				log.Println(
					"PROCESS: ",
					r.ID,
					r.UserID,
					r.ProductID,
					r.ProductVariantID,
				)
				emailData.Status = "✅ Aktif"
			}

			err = model.UpdateLastSentAt(r.ID, time.Now())
			if err != nil {
				defer errorUpdateReminderSchedules()
				errMessage := "Failed to update data in table reminder_schedules: " + err.Error()
				panic(errMessage)
			}
			htmlBody, err := RenderSubscriptionEmail(emailData, os.Getenv("TEMPLATES_PATH")+"/email_schedule_reminder.html")
			if err != nil {
				defer errorSubscriptionReminderEmail()
				errMessage := "Failed to parse html file for subscription reminder email template: " + err.Error()
				panic(errMessage)
			}

			sender := NewBrevoSender()
			_ = sender.SendMail(
				r.UEmail,
				"Pengingat Langganan Product Subcommerce",
				htmlBody,
			)
		}
	}
}

func checkSchedule(t time.Time, today string) bool {
	return t.Format("2006-01-02") == today
}
