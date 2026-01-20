package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/errs"
)

func GetAllSubscriptionsByUser(userId interface{}) ([]entity.UserSubscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var subs []entity.UserSubscription
	err := db.DB.SelectContext(ctx, &subs, `SELECT 
	rs.id, rs.product_id, rs.product_variant_id, rs.last_sent_at, rs.next_send, rs.next_warning_send, rs.next_remove, rs.is_over,
	o.order_pretty_id, pay.name AS payment_name, os.name AS order_status_name,
	s.name AS seller_name, s.img as seller_img,
	p.name AS product_name, pv.name AS product_variant_name, pi.img AS product_img, 
	o.quantity, o.note, pv.price, pv.discount, pv.interval, i.name AS interval_name
		FROM reminder_schedules rs
		JOIN orders o ON rs.id = o.order_uq_id
		JOIN order_statuses os ON o.order_status_id = os.id
		JOIN payments pay ON o.payment_id = pay.id
		JOIN products p ON o.product_id = p.id
		JOIN product_variants pv ON o.product_variant_id = pv.id
		JOIN LATERAL (SELECT pi.img FROM product_images pi WHERE p.id = pi.product_id LIMIT 1) pi ON true
		JOIN sellers s ON p.seller_id = s.id
		JOIN intervals i ON pv.interval_id = i.id
		WHERE rs.user_id = $1
		ORDER BY rs.created_at DESC, o.order_pretty_id DESC;`, userId)
	if err != nil {
		return nil, err
	}

	if len(subs) == 0 {
		return nil, errs.ErrNoSubscriptionFound
	}

	return subs, nil
}

func DeleteReminderSchedule(orderId string) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.DB.QueryRowContext(ctx, `DELETE FROM reminder_schedules
	WHERE order_uq_id = $2
	RETURNING order_uq_id`, orderId).Scan(&orderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoOrderFound
		}
		return nil, err
	}

	return &orderId, nil
}
