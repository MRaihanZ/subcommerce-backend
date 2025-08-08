package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sqlx.DB

func init() {
	var err error

	// Load .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		errMessage := "Failed to connect to database: " + err.Error()
		log.Fatalln(errMessage)
	}

	log.Println("Database connected")
}

func users() {
	fmt.Println("")
	log.Println("+++ Seeding users table +++")
	gofakeit.Seed(0)
	var id uuid.UUID
	// Generate DOB for age between 18 and 60
	maxDOB := time.Now().AddDate(-18, 0, 0) // 18 years ago
	minDOB := time.Now().AddDate(-60, 0, 0) // 60 years ago
	var dob time.Time

	password := "123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		id = uuid.New()
		dob = gofakeit.DateRange(minDOB, maxDOB)
		_, err := DB.Exec(`INSERT INTO users (id, name, email, dob, password)
		VALUES ($1, $2, $3, $4, $5)`,
			id, gofakeit.Name(), gofakeit.Email(), dob, hashedPassword)
		if err != nil {
			log.Println("insert error: ", err)
		}
	}
	log.Println("=== Complete ===")
}

func sellers() {
	fmt.Println("")
	log.Println("+++ Seeding sellers table +++")
	gofakeit.Seed(0)

	var id uuid.UUID

	var users []string
	err := DB.Select(&users, "SELECT id FROM users LIMIT 5")
	if err != nil {
		log.Fatalf("Failed to fetch sellers: %v", err)
	}

	sold_products := []int{234, 167, 180, 199, 241}
	average_rating := []float32{4.799, 3.872, 3.117, 4.929, 4.243}
	rating_total := []int{190, 50, 70, 445, 290}
	rating_count := []int{100, 32, 30, 89, 59}
	for i, v := range users {
		id = uuid.New()
		_, err := DB.Exec(`INSERT INTO sellers (id, user_id, name, address, sold_products, average_rating, rating_total, rating_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			id, v, gofakeit.Name(), gofakeit.Address().Address, sold_products[i], average_rating[i], rating_total[i], rating_count[i])
		if err != nil {
			log.Println("insert error: ", err)
		}
	}
	log.Println("=== Complete ===")
}

func admins() {
	fmt.Println("")
	log.Println("+++ Seeding admin table +++")

	password := "123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	id := uuid.New()
	_, err = DB.Exec(`INSERT INTO admins (id, name, password)
		VALUES ($1, $2, $3)`,
		id, gofakeit.Name(), hashedPassword)
	if err != nil {
		log.Println("insert error: ", err)
	}
	log.Println("=== Complete ===")
}

func intervals() {
	fmt.Println("")
	log.Println("+++ Seeding intervals table +++")
	code := []string{"D", "W", "M", "Y"}
	name := []string{"hari", "minggu", "bulan", "tahun"}

	for i, v := range code {
		_, err := DB.Exec(`INSERT INTO intervals (code, name)
		VALUES ($1, $2)`,
			v, name[i])
		if err != nil {
			log.Println("insert error: ", err)
		}
	}
	log.Println("=== Complete ===")
}

func products() {
	fmt.Println("")
	log.Println("+++ Seeding product table +++")

	gofakeit.Seed(0)

	var sellers []string
	err := DB.Select(&sellers, "SELECT id FROM sellers LIMIT 5")
	if err != nil {
		log.Fatalf("Failed to fetch sellers: %v", err)
	}

	sold := []int{100, 32, 34, 89, 61}
	average_rating := []float32{4.599, 3.456, 2.604, 5.00, 4.961}
	rating_total := []int{190, 50, 70, 445, 290}
	rating_count := []int{100, 32, 30, 89, 59}

	for i, v := range sellers {
		_, err := DB.Exec(`INSERT INTO products (seller_id, name, description, sold, average_rating, rating_total, rating_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			v, gofakeit.ProductName(), gofakeit.ProductDescription(), sold[i], average_rating[i], rating_total[i], rating_count[i])
		if err != nil {
			log.Println("insert error: ", err)
		}
	}
	log.Println("=== Complete ===")
}

func product_variants() {
	fmt.Println("")
	log.Println("+++ Seeding product_variants table +++")

	gofakeit.Seed(0)

	var products []string
	err := DB.Select(&products, "SELECT id FROM products LIMIT 3")
	if err != nil {
		log.Fatalf("Failed to fetch products: %v", err)
	}

	variants := []string{"A", "B"}

	for _, p := range products {
		for _, v := range variants {
			_, err := DB.Exec(`INSERT INTO product_variants (product_id, interval_id, is_default, name, interval, stock, sold, price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
				p, 3, false, v, 1, 33, 11, 120000)
			if err != nil {
				log.Println("insert error: ", err)
			}
		}
	}

	def := []int{4, 5}

	for _, d := range def {
		_, err := DB.Exec(`INSERT INTO product_variants (product_id, interval_id, is_default, name, interval, stock, sold, price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			d, 3, true, "default", 1, 100, 200, 240000)
		if err != nil {
			log.Println("insert error: ", err)
		}
	}
	log.Println("=== Complete ===")
}

func product_images() {
	fmt.Println("")
	log.Println("+++ Seeding product_images table +++")

	gofakeit.Seed(0)

	var products []string
	err := DB.Select(&products, "SELECT id FROM products LIMIT 5")
	if err != nil {
		log.Fatalf("Failed to fetch products: %v", err)
	}

	for _, p := range products {
		for i := 0; i < 3; i++ {
			_, err := DB.Exec(`INSERT INTO product_images (product_id, img)
		VALUES ($1, $2)`,
				p, "/assets/img/item.jpg")
			if err != nil {
				log.Println("insert error: ", err)
			}
		}
	}
	log.Println("=== Complete ===")
}

func ratings() {
	fmt.Println("")
	log.Println("+++ Seeding ratings table +++")

	gofakeit.Seed(0)

	var products []string
	err := DB.Select(&products, "SELECT id FROM products LIMIT 5")
	if err != nil {
		log.Fatalf("Failed to fetch products: %v", err)
	}

	var product_variants []string
	err = DB.Select(&product_variants, "SELECT id FROM product_variants")
	if err != nil {
		log.Fatalf("Failed to fetch product_variants: %v", err)
	}

	var users []string
	err = DB.Select(&users, "SELECT id FROM users")
	if err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}

	ratings := []int{1, 2, 3, 4, 5}
	rating := gofakeit.RandomInt(ratings)

	for i, p := range products {
		for _, pv := range product_variants {
			_, err := DB.Exec(`INSERT INTO ratings (product_id, product_variant_id, user_id, rating, comment)
		VALUES ($1, $2, $3, $4, $5)`,
				p, pv, users[i], rating, gofakeit.Comment())
			if err != nil {
				log.Println("insert error: ", err)
			}
		}
	}
	log.Println("=== Complete ===")
}

func checkout_statuses() {
	fmt.Println("")
	log.Println("+++ Seeding checkout statuses table +++")
	_, err := DB.Exec(`INSERT INTO checkout_statuses (name)
		VALUES ('menunggu pembayaran'),
		('pembayaran dibatalkan'),
		('batas waktu pembayaran habis'),
		('menunggu konfirmasi seller'),
		('dibatalkan seller'),
		('dibatalkan pengguna'),
		('produk sedang disiapkan'),
		('produk sudah dikirim'),
		('produk tidak diterima'),
		('pesanan selesai')`)
	if err != nil {
		log.Println("insert error: ", err)
	}
	log.Println("=== Complete ===")
}

func payments() {
	fmt.Println("")
	log.Println("+++ Seeding checkout statuses table +++")
	_, err := DB.Exec(`INSERT INTO checkout_statuses (name)
		VALUES 
		('Qris', '
	/assets/svg/payments_24dp_E3E3E3_FILL0_wght400_GRAD0_opsz24.svg'),
		('Transfer Mandiri', '
	/assets/svg/payments_24dp_E3E3E3_FILL0_wght400_GRAD0_opsz24.svg'),
		('pembayaran dibatalkan', '
	/assets/svg/payments_24dp_E3E3E3_FILL0_wght400_GRAD0_opsz24.svg')`)
	if err != nil {
		log.Println("insert error: ", err)
	}
	log.Println("=== Complete ===")
}

func main() {
	users()
	sellers()
	admins()
	intervals()
	products()
	product_variants()
	product_images()
	ratings()
	checkout_statuses()
}
