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
		JOIN products p ON rs.product_id = p.id
		JOIN product_variants pv ON rs.product_variant_id = pv.id
		JOIN LATERAL (SELECT pi.img FROM product_images pi WHERE p.id = pi.product_id LIMIT 1) pi ON true
		JOIN sellers s ON p.seller_id = s.id
		JOIN intervals i ON pv.interval_id = i.id
		WHERE rs.user_id = $1;`, userId)
	if err != nil {
		return nil, err
	}

	if len(subs) == 0 {
		return nil, errs.ErrNoSubscriptionFound
	}

	return subs, nil
}

func GetAllSubscriptionsBySeller(sellerId interface{}) ([]entity.UserSubscriptionBySeller, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var subs []entity.UserSubscriptionBySeller
	err := db.DB.SelectContext(ctx, &subs, `SELECT 
	rs.id, rs.product_id, rs.product_variant_id, rs.last_sent_at, rs.next_send, rs.next_warning_send, rs.next_remove, rs.is_over,
	o.order_pretty_id, pay.name AS payment_name, os.name AS order_status_name,
	u.name AS user_name, u.img as user_img,
	p.name AS product_name, pv.name AS product_variant_name, pi.img AS product_img, 
	o.quantity, o.note, pv.price, pv.discount, pv.interval, i.name AS interval_name
		FROM reminder_schedules rs
		JOIN orders o ON rs.id = o.order_uq_id
		JOIN order_statuses os ON o.order_status_id = os.id
		JOIN payments pay ON o.payment_id = pay.id
		JOIN products p ON rs.product_id = p.id
		JOIN product_variants pv ON rs.product_variant_id = pv.id
		JOIN LATERAL (SELECT pi.img FROM product_images pi WHERE p.id = pi.product_id LIMIT 1) pi ON true
		JOIN sellers s ON p.seller_id = s.id
    JOIN users u ON rs.user_id = u.id
		JOIN intervals i ON pv.interval_id = i.id
		WHERE s.id = $1;`, sellerId)
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

	var returnId string
	err := db.DB.QueryRowContext(ctx, `DELETE FROM reminder_schedules
	WHERE id = $1
	RETURNING id`, orderId).Scan(&returnId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoSubscriptionFound
		}
		return nil, err
	}

	return &returnId, nil
}

func GetEmailSellerByOrderId(orderId string) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var email string
	err := db.DB.GetContext(ctx, &email, `SELECT u.email FROM orders o
	JOIN products p ON o.product_id = p.id
	JOIN sellers s ON p.seller_id = s.id
	JOIN users u ON s.user_id = u.id
	WHERE o.order_uq_id = $1;`, orderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoSubscriptionFound
		}
		return nil, err
	}

	return &email, nil
}

func GetReminderScheduleId(userId interface{}) (*entity.SubscriptionIdPaymentLink, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var data entity.SubscriptionIdPaymentLink
	err := db.DB.GetContext(ctx, &data, `SELECT rs.id, o.payment_link FROM reminder_schedules rs
	JOIN orders o ON o.order_uq_id = rs.id
	WHERE rs.user_id = $1 ORDER BY rs.created_at DESC;`, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoSubscriptionFound
		}
		return nil, err
	}

	return &data, nil
}

func UpdateReminderScheduleInterval(id string, next, warning, remove time.Time) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnId string
	err := db.DB.QueryRowContext(ctx, `UPDATE reminder_schedules 
	SET next_send = $1, next_warning_send = $2,
	next_remove = $3
	WHERE id = $4
	RETURNING id`, next, warning, remove, id).Scan(&returnId)
	if err != nil {
		return nil, err
	}

	return &returnId, nil
}
