package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/errs"
)

func GetAllOrder(userId interface{}) ([]entity.OrderGetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var orders []entity.OrderGetResponse
	err := db.DB.SelectContext(ctx, &orders, `SELECT pay.name AS pay_name, o.id AS order_id, os.name AS os_name,
	o.rating, o.product_id, o.product_variant_id, o.quantity, o.total_price, o.order_pretty_id, o.note, o.payment_link, o.order_uq_id,
	s.name AS s_name, s.img as s_img, p.name AS p_name, pi.img AS p_img, p.active,
	pv.name AS pv_name, pv.interval, i.name AS i_name
	FROM orders o
	JOIN order_statuses os ON o.order_status_id = os.id
	JOIN payments pay ON o.payment_id = pay.id
	JOIN products p ON o.product_id = p.id
	JOIN product_variants pv ON o.product_variant_id = pv.id
	JOIN LATERAL (SELECT pi.img FROM product_images pi WHERE p.id = pi.product_id LIMIT 1) pi ON true
	JOIN sellers s ON p.seller_id = s.id
	JOIN intervals i ON pv.interval_id = i.id
	WHERE o.user_id = $1
	ORDER BY o.created_at DESC, o.order_pretty_id DESC;`, userId)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, errs.ErrNoOrderFound
	}

	return orders, nil
}

func GetOrderBySellerID(sellerId interface{}, stateAction, orderCreated, orderId string) ([]entity.OrderGetResponseBySeller, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var orders []entity.OrderGetResponseBySeller
	var query string
	var args []interface{}

	switch stateAction {
	case "next":
		// Get the next order
		query = `SELECT o.id AS order_id, o.order_pretty_id, u.name AS u_name, u.img AS u_img, o.product_id, 
		o.product_variant_id, p.name AS p_name, pv.name AS pv_name, pi.img AS p_img, o.quantity,
		pv.interval, i.name AS i_name, pay.name AS pay_name, os.name AS os_name, o.total_price, o.created_at
		FROM orders o
		JOIN order_statuses os ON o.order_status_id = os.id
		JOIN payments pay ON o.payment_id = pay.id
		JOIN products p ON o.product_id = p.id
		JOIN product_variants pv ON o.product_variant_id = pv.id
		JOIN LATERAL (SELECT pi.img FROM product_images pi WHERE p.id = pi.product_id LIMIT 1) pi ON true
		JOIN sellers s ON p.seller_id = s.id
		JOIN users u ON o.user_id = u.id
		JOIN intervals i ON pv.interval_id = i.id
		WHERE s.id = $1
		AND (o.created_at, o.id) <= ($2, $3)
		ORDER BY o.created_at DESC, o.order_pretty_id DESC
		LIMIT 10;`
		args = []interface{}{sellerId, orderCreated, orderId}
	case "previous":
		// Get the previous order
		query = `SELECT o.id AS order_id, o.order_pretty_id, u.name AS u_name, u.img AS u_img, o.product_id, 
		o.product_variant_id, p.name AS p_name, pv.name AS pv_name, pi.img AS p_img, o.quantity,
		pv.interval, i.name AS i_name, pay.name AS pay_name, os.name AS os_name, o.total_price, o.created_at
		FROM orders o
		JOIN order_statuses os ON o.order_status_id = os.id
		JOIN payments pay ON o.payment_id = pay.id
		JOIN products p ON o.product_id = p.id
		JOIN product_variants pv ON o.product_variant_id = pv.id
		JOIN LATERAL (SELECT pi.img FROM product_images pi WHERE p.id = pi.product_id LIMIT 1) pi ON true
		JOIN sellers s ON p.seller_id = s.id
		JOIN users u ON o.user_id = u.id
		JOIN intervals i ON pv.interval_id = i.id
		WHERE s.id = $1
		AND (o.created_at, o.id) >= ($2, $3)
		ORDER BY o.created_at ASC, o.order_pretty_id ASC
		LIMIT 10;`
		args = []interface{}{sellerId, orderCreated, orderId}
	default:
		//get the first order
		query = `SELECT o.id AS order_id, o.order_pretty_id, u.name AS u_name, u.img AS u_img, o.product_id, 
		o.product_variant_id, p.name AS p_name, pv.name AS pv_name, pi.img AS p_img, o.quantity,
		pv.interval, i.name AS i_name, pay.name AS pay_name, os.name AS os_name, o.total_price, o.created_at
		FROM orders o
		JOIN order_statuses os ON o.order_status_id = os.id
		JOIN payments pay ON o.payment_id = pay.id
		JOIN products p ON o.product_id = p.id
		JOIN product_variants pv ON o.product_variant_id = pv.id
		JOIN LATERAL (SELECT pi.img FROM product_images pi WHERE p.id = pi.product_id LIMIT 1) pi ON true
		JOIN sellers s ON p.seller_id = s.id
		JOIN users u ON o.user_id = u.id
		JOIN intervals i ON pv.interval_id = i.id
		WHERE s.id = $1
		ORDER BY o.created_at DESC, o.order_pretty_id DESC
		LIMIT 10;`
		args = []interface{}{sellerId}
	}

	err := db.DB.SelectContext(ctx, &orders, query, args...)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, errs.ErrNoOrderFound
	}

	return orders, nil
}

func CreateOrder(id interface{}, order []entity.OrderRequest, paymentUrl string, orderId interface{}, state string) (*string, error) {
	query := "INSERT INTO orders (user_id, payment_id, order_status_id, product_id, product_variant_id, payment_link, note, quantity, unit_price, total_price, order_uq_id, order_pretty_id) VALUES "
	args := []interface{}{}
	reqData := []string{}

	for i, arg := range order {
		n := i*11 + 1
		reqData = append(reqData, fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, 'INV-' || TO_CHAR(NOW(), 'YYYYMMDD') || '-' ||
  		LPAD(nextval('orders_order_pretty_id_seq')::text, 4, '0'))`, n, n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8, n+9, n+10))
		args = append(args, id, arg.PayId, 1, arg.PId, arg.PVId, paymentUrl, arg.Note, arg.Quantity, arg.UnitPrice, arg.TotalPrice, orderId)
	}
	query += strings.Join(reqData, ",")
	query += "RETURNING user_id"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userId string

	err := db.DB.QueryRowxContext(ctx, query, args...).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoOrderFound
		}
		return nil, err
	}

	if state == "order" {
		_, err = DeleteCheckoutOrder(id)
		if err != nil {
			if errors.Is(err, errs.ErrNoCheckoutFound) {
				return nil, errs.ErrNoCheckoutFound
			}
			return nil, err
		}
	}

	return &userId, nil
}

func DeleteOrder(id interface{}) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnId string

	err := db.DB.QueryRowContext(ctx, `DELETE FROM checkouts WHERE user_id = $1
	RETURNING user_id`, id).Scan(&returnId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoOrderFound
		}
		return nil, err
	}

	return &returnId, nil
}

func UpdateStatusOrderBySeller(orderId string, statusId int) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.DB.QueryRowContext(ctx, `UPDATE orders SET order_status_id = $1
	WHERE id = $2
	RETURNING order_uq_id`, statusId, orderId).Scan(&orderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoOrderFound
		}
		return nil, err
	}

	return &orderId, nil
}

func GetUserIdFromOrders(orderUgId string) (*entity.GetUserProductProductVariant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var res entity.GetUserProductProductVariant
	err := db.DB.GetContext(ctx, &res, `SELECT user_id, product_id, product_variant_id
	FROM orders WHERE order_uq_id = $1`, orderUgId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoOrderFound
		}
		return nil, err
	}

	return &res, nil
}

func GetIntervalProduct(productId, productVariantId int) (*entity.IntervalProduct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var res entity.IntervalProduct
	err := db.DB.GetContext(ctx, &res, `SELECT interval_id, interval
	FROM product_variants WHERE product_id = $1 AND id = $2`, productId, productVariantId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoProductFound
		}
		return nil, err
	}

	return &res, nil
}

func CreateReminderSchedule(orderUqId string, data entity.GetUserProductProductVariant, next, warning, remove time.Time) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var reminderSchedulesId string

	err := db.DB.QueryRowxContext(ctx, "INSERT INTO reminder_schedules (id, user_id, product_id, product_variant_id, next_send, next_warning_send, next_remove) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		orderUqId, data.UserId, data.ProductId, data.ProductVariantId, next, warning, remove,
	).Scan(&reminderSchedulesId)
	if err != nil {
		return nil, err
	}
	return &reminderSchedulesId, nil
}

func UpdateStatusOrder(orderId string, statusId int) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.DB.QueryRowContext(ctx, `UPDATE orders SET order_status_id = $1
	WHERE order_uq_id = $2
	RETURNING order_status_id`, statusId, orderId).Scan(&orderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoOrderFound
		}
		return nil, err
	}

	return &orderId, nil
}

func UpdateRatingOrder(productId int, productVariantId int, orderId int, userId interface{}) (*bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rating bool
	err := db.DB.QueryRowContext(ctx, `UPDATE orders SET rating = $1
	WHERE product_id = $2 AND product_variant_id = $3 AND id = $4 AND user_id = $5
	RETURNING rating`, true, productId, productVariantId, orderId, userId).Scan(&rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoOrderFound
		}
		return nil, err
	}

	return &rating, nil
}

func GetAllCheckoutOrder(userId interface{}) ([]entity.GetCheckoutOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var checkouts []entity.GetCheckoutOrderResponse
	err := db.DB.SelectContext(ctx, &checkouts, `SELECT c.product_id, c.product_variant_id, s.name AS s_name, s.img AS s_img,
	p.name AS p_name, pi.img AS p_img, pv.name AS pv_name, pv.interval, i.name AS i_name,
	c.quantity, c.total_price
	FROM checkouts c
	JOIN products p ON c.product_id = p.id
	JOIN LATERAL (SELECT pi.img FROM product_images pi WHERE c.product_id = pi.product_id LIMIT 1) pi ON true
	JOIN product_variants pv ON c.product_variant_id = pv.id
	JOIN sellers s ON p.seller_id = s.id
	JOIN intervals i ON pv.interval_id = i.id
	WHERE c.user_id = $1 AND p.active = true
	ORDER BY c.created_at DESC;`, userId)
	if err != nil {
		return nil, err
	}

	if len(checkouts) == 0 {
		return nil, errs.ErrNoCheckoutFound
	}

	return checkouts, nil
}

func GetAllCheckoutOrderPayment(userId interface{}) ([]entity.GetCheckoutOrderPaymentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var checkouts []entity.GetCheckoutOrderPaymentResponse
	err := db.DB.SelectContext(ctx, &checkouts, `SELECT c.product_id, c.product_variant_id, s.name AS s_name, s.img AS s_img,
	p.name AS p_name, pi.img AS p_img, pv.name AS pv_name, pv.interval, i.name AS i_name, pv.discount,
	c.quantity, c.total_price, c.unit_price
	FROM checkouts c
	JOIN products p ON c.product_id = p.id
	JOIN LATERAL (SELECT pi.img FROM product_images pi WHERE c.product_id = pi.product_id LIMIT 1) pi ON true
	JOIN product_variants pv ON c.product_variant_id = pv.id
	JOIN sellers s ON p.seller_id = s.id
	JOIN intervals i ON pv.interval_id = i.id
	WHERE c.user_id = $1 AND p.active = true
	ORDER BY c.created_at DESC;`, userId)
	if err != nil {
		return nil, err
	}

	if len(checkouts) == 0 {
		return nil, errs.ErrNoCheckoutFound
	}

	return checkouts, nil
}

func DeleteCheckoutOrder(userId interface{}) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var DeleteCheckoutUser string
	err := db.DB.QueryRowxContext(ctx, `DELETE FROM checkouts
	WHERE user_id = $1
	RETURNING user_id`, userId).Scan(&DeleteCheckoutUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoCheckoutFound
		}
		return nil, err
	}

	return &DeleteCheckoutUser, nil
}

func CreateCheckoutOrder(id interface{}, checkouts []entity.OrderCheckout, stateAction string) (*string, error) {
	if len(checkouts) == 0 {
		return nil, errs.ErrCheckoutRequestZero
	}

	_, err := GetAllCheckoutOrder(id)
	if err != nil {
		if errors.Is(err, errs.ErrNoCheckoutFound) {
			query := "INSERT INTO checkouts (user_id, product_id, product_variant_id, quantity, unit_price, total_price) VALUES "
			args := []interface{}{}
			reqData := []string{}

			for i, arg := range checkouts {
				n := i*6 + 1
				reqData = append(reqData, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", n, n+1, n+2, n+3, n+4, n+5))
				args = append(args, id, arg.PId, arg.PVId, arg.Quantity, arg.UnitPrice, arg.TotalPrice)
			}
			query += strings.Join(reqData, ",")
			query += "RETURNING user_id"

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var userId string

			err = db.DB.QueryRowxContext(ctx, query, args...).Scan(&userId)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, errs.ErrNoCheckoutFound
				}
				return nil, err
			}

			if stateAction == "cart" {
				for _, val := range checkouts {
					_, err = DeleteCart(id, val.PId, val.PVId)
					if err != nil {
						return nil, err
					}
				}
			}
			return &userId, nil
		}
		return nil, err
	}

	_, err = DeleteCheckoutOrder(id)
	if err != nil {
		if errors.Is(err, errs.ErrNoCheckoutFound) {
			return nil, errs.ErrNoCheckoutFound
		}
		return nil, err
	}

	query := "INSERT INTO checkouts (user_id, product_id, product_variant_id, quantity, unit_price, total_price) VALUES "
	args := []interface{}{}
	reqData := []string{}

	for i, arg := range checkouts {
		n := i*6 + 1
		reqData = append(reqData, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", n, n+1, n+2, n+3, n+4, n+5))
		args = append(args, id, arg.PId, arg.PVId, arg.Quantity, arg.UnitPrice, arg.TotalPrice)
	}
	query += strings.Join(reqData, ",")
	query += "RETURNING user_id"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userId string

	err = db.DB.QueryRowxContext(ctx, query, args...).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoCheckoutFound
		}
		return nil, err
	}
	if stateAction == "cart" {
		for _, val := range checkouts {
			_, err = DeleteCart(id, val.PId, val.PVId)
			if err != nil {
				return nil, err
			}
		}
	}
	return &userId, nil
}

func GetAllOrderPayment() ([]entity.GetOrderPaymentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var payments []entity.GetOrderPaymentResponse
	err := db.DB.SelectContext(ctx, &payments, `SELECT p.id AS p_id, p.name AS p_name, p.category_payment_id, cp.name AS cp_name
	FROM payments p
	JOIN category_payments cp ON p.category_payment_id = cp.id
	ORDER BY p.category_payment_id;`)
	if err != nil {
		return nil, err
	}

	if len(payments) == 0 {
		return nil, errs.ErrNoPaymentFound
	}

	return payments, nil
}

func GetExistingOrderSubscriptionByUserId(userId interface{}, productId, productVariantId int) (*bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result bool
	err := db.DB.GetContext(ctx, &result, `SELECT EXISTS (
    SELECT 1
    FROM orders
    WHERE (order_status_id = 13 OR order_status_id = 10) AND user_id = $1 AND product_id = $2 AND product_variant_id = $3
	LIMIT 1
	)`, userId, productId, productVariantId)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetExistingOrderSubscription(orderId interface{}) (*bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var res entity.GetUserProductProductVariant
	err := db.DB.GetContext(ctx, &res, `SELECT user_id, product_id, product_variant_id
	FROM orders WHERE order_uq_id = $1`, orderId)
	if err != nil {
		return nil, err
	}

	var result bool
	err = db.DB.GetContext(ctx, &result, `SELECT EXISTS (
    SELECT 1
    FROM orders
    WHERE (order_status_id = 13 OR order_status_id = 10) AND user_id = $1 AND product_id = $2 AND product_variant_id = $3
	LIMIT 1
	)`, res.UserId, res.ProductId, res.ProductVariantId)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateReminderScheduleId(orderId interface{}) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get user_id, product_id, product_variant_id from recent order
	var res entity.GetUserProductProductVariant
	err := db.DB.GetContext(ctx, &res, `SELECT user_id, product_id, product_variant_id
	FROM orders WHERE order_uq_id = $1`, orderId)
	if err != nil {
		return nil, err
	}

	// get the latest order_uq_id before recent order
	var result []string
	err = db.DB.SelectContext(ctx, &result, `SELECT order_uq_id
	FROM orders
	WHERE user_id = $1 AND product_id = $2 AND product_variant_id = $3
	ORDER BY created_at DESC
	LIMIT 2`, res.UserId, res.ProductId, res.ProductVariantId)
	if err != nil {
		return nil, err
	}

	var newOrderId string
	err = db.DB.QueryRowContext(ctx, `UPDATE reminder_schedules SET id = $1
	WHERE id = $2 RETURNING id;`, orderId, result[1]).Scan(&newOrderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoOrderFound
		}
		return nil, err
	}

	return &newOrderId, nil
}

func GetNextRemoveReminderSchedule(orderUqId string) (*time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var nextRemove time.Time
	err := db.DB.GetContext(ctx, &nextRemove, `SELECT next_remove
	FROM reminder_schedules
	WHERE id = $1;`, orderUqId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoSubscriptionFound
		}
		return nil, err
	}

	return &nextRemove, nil
}

func UpdateWalletSeller(orderId string) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var totalPrice int64
	err := db.DB.GetContext(ctx, &totalPrice, `SELECT total_price
	FROM orders WHERE id = $1`, orderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoOrderFound
		}
		return nil, err
	}

	var sellerId string
	err = db.DB.QueryRowContext(ctx, `UPDATE sellers s
	SET wallet = $1
	FROM orders o
	JOIN products p ON p.id = o.product_id
	WHERE s.id = p.seller_id
	AND o.id = $2
	RETURNING s.id;`, totalPrice, orderId).Scan(&sellerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoOrderFound
		}
		return nil, err
	}

	return &sellerId, nil
}

// func UpdateReminderScheduleSendEmail(orderId interface{}) (*string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	// get user_id, product_id, product_variant_id from recent order
// 	var res entity.GetUserProductProductVariant
// 	err := db.DB.GetContext(ctx, &res, `SELECT user_id, product_id, product_variant_id
// 	FROM orders WHERE order_uq_id = $1`, orderId)
// 	if err != nil {
// 		return nil, err
// 	}
// }
