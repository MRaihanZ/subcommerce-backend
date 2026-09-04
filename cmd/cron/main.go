package main

import (
	"log"
	"os"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/service"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func init() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
}

func main() {
	log.Println("✅ Starting Cron Job")
	loc, _ := time.LoadLocation("Asia/Jakarta")

	// if second enable total expression 6 if not 5
	// expression will be 6 or 5 *: (* * * * * *) or (* * * * *)
	// if second enable
	// 1 * (0-59) = in second
	// if second disable
	// 2 * (0-59) = in minute
	// 3 * (0-23) = in hour
	// 4 * (1-31) = in day
	// 5 * (1-12) = in month
	// 6 or 5 * = weekday
	// not every second or expression but "at" expression
	// example "(* 20 * * * *)"
	// will run every hour at minute 20, example: 07:20:00 or 19:20:00
	c := cron.New(
		cron.WithLocation(loc),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
		),

		// second enable
		// cron.WithSeconds(),
	)

	// Run every day at 07:00
	// Run every second
	// To run every one minute "*/1 * * * *"
	// To run every day at 07:00 "0 7 * * *"
	c.AddFunc("0 7 * * *", func() {
		log.Println("✅ Starting Reminder Cron")
		service.RunReminderCron()
		log.Println("✅ Finish Reminder Cron")
	})

	c.Start()
	select {}
}
