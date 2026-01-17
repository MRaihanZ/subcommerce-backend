package main

import (
	"log"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/service"
)

func main() {
	// Status Subscription
	// 1. ✅ Aktif
	// 2. ⚠️ Akan Berakhir
	// 3. ❌ Tidak Aktif
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
		Domain:          "subcommerce.mraihanz.my.id",
	}

	htmlBody, err := service.RenderSubscriptionEmail(emailData, "E:/GIU/Devel/go_app/subcommerce-backend/internal/templates/email_schedule_reminder.html")
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
