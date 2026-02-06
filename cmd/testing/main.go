package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/service"
	"github.com/MRaihanZ/subcommerce-backend/internal/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func SendMail() {
	// Status Subscription
	// 1. ✅ Aktif
	// 2. ⚠️ Akan Berakhir
	// 3. ❌ Tidak Aktif
	today := time.Now()
	emailData := entity.SubscriptionEmailData{
		OrderID:         "INV-20260121-0049",
		SubscriptionID:  "557ba9c5-c613-4131-aa69-140f72b03877",
		UserName:        "Raihan",
		ProductName1:    "Iron Nexus Plastic",
		ProductName2:    "A",
		Subscription:    "2 Bulan",
		Status:          "✅ Aktif",
		OrderDate:       today.Format("02-01-2006 15:04:05"),
		PaymentDeadline: "2026-01-19",
		TargetEmail:     "mraihanzhafran.14@gmail.com",
		Domain:          os.Getenv("WEBSITE_URL"),
	}

	htmlBody, err := service.RenderSubscriptionEmail(emailData, os.Getenv("TEMPLATES_PATH")+"/email_schedule_reminder.html")
	if err != nil {
		log.Fatal(err)
	}

	sender := service.NewBrevoSender()
	_ = sender.SendMail(
		"mraihanzhafran.14@gmail.com",
		"Pengingat Langganan Product Subcommerce",
		htmlBody,
	)
}

func IntervalSet() {
	start := time.Now()

	next, warning, remove := utils.CalculateReminderDates(
		start,
		1,
		8,
	)

	log.Println("Next: ", next)
	log.Println("Warning: ", warning)
	log.Println("Remove: ", remove)

	// store next, warning, remove to DB

}

func CasesTester() {
	caser := cases.Title(language.MustParse("id-ID"))
	some := "Bulan"
	fmt.Println(caser.String(some))
}

func main() {
	// SendMail()
	// IntervalSet()
	CasesTester()
}
