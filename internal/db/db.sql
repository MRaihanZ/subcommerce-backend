CREATE DATABASE subcommerce;

use subcommerce;

-- MAKE SURE TO CHANGE IMG DEFAULT PATH
CREATE TABLE
    users (
        id UUID PRIMARY KEY,
        name VARCHAR(50) NOT NULL,
        img TEXT DEFAULT '/assets/img/profile2.jpg' NOT NULL,
        email TEXT NOT NULL,
        dob DATE NOT NULL,
        -- dob = date of birth
        password VARCHAR(60) NOT NULL,
        created_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

-- MAKE SURE TO CHANGE IMG DEFAULT PATH
CREATE TABLE
    sellers (
        id UUID PRIMARY KEY,
        user_id UUID CONSTRAINT fk_s_user_id REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        CONSTRAINT uq_s_seller_user UNIQUE (id, user_id),
        name VARCHAR(50) NOT NULL,
        img TEXT DEFAULT '/assets/img/profile1.jpg' NOT NULL,
        address TEXT NOT NULL,
        sold_products INT DEFAULT 0 NOT NULL,
        average_rating NUMERIC(4, 3) DEFAULT 0 NOT NULL,
        rating_total INT DEFAULT 0 NOT NULL,
        rating_count INT DEFAULT 0 NOT NULL,
        created_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

CREATE TABLE
    admins (
        id UUID PRIMARY KEY,
        name VARCHAR(50) NOT NULL,
        password VARCHAR(60) NOT NULL
    );

CREATE TABLE
    intervals (
        id SERIAL PRIMARY KEY,
        code VARCHAR(1),
        name VARCHAR(6)
    );

CREATE TABLE
    products (
        id SERIAL PRIMARY KEY,
        seller_id UUID CONSTRAINT fk_p_seller_id REFERENCES sellers (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        name VARCHAR(100) NOT NULL,
        description TEXT NOT NULL,
        sold INT DEFAULT 0 NOT NULL,
        average_rating NUMERIC(4, 3) DEFAULT 0 NOT NULL,
        rating_total INT DEFAULT 0 NOT NULL,
        rating_count INT DEFAULT 0 NOT NULL,
        active BOOLEAN DEFAULT TRUE NOT NULL,
        created_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

CREATE TABLE
    product_variants (
        id SERIAL PRIMARY KEY,
        product_id INT CONSTRAINT fk_pv_product_id REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        interval_id INT CONSTRAINT fk_pv_interval_id REFERENCES intervals (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        is_default BOOLEAN DEFAULT TRUE NOT NULL, -- If true, the variant will be hidden. if false, it will be shown.
        name VARCHAR(20),
        interval SMALLINT,
        stock SMALLINT DEFAULT 0 NOT NULL,
        sold INT DEFAULT 0 NOT NULL,
        price BIGINT DEFAULT 0 NOT NULL,
        discount SMALLINT DEFAULT 0 NOT NULL,
        min_order SMALLINT DEFAULT 1 NOT NULL,
        created_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

-- MAKE SURE TO CHANGE IMG DEFAULT PATH
CREATE TABLE
    product_images (
        id SERIAL PRIMARY KEY,
        product_id INT CONSTRAINT fk_pi_product_id REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        img TEXT NOT NULL
    );

CREATE TABLE
    ratings (
        id SERIAL PRIMARY KEY,
        product_id INT CONSTRAINT fk_ra_product_id REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        product_variant_id INT CONSTRAINT fk_ra_product_variant_id REFERENCES product_variants (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        user_id UUID CONSTRAINT fk_ra_user_id REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        CONSTRAINT uq_ra_product_product_variant_user UNIQUE (product_id, product_variant_id, user_id),
        rating INT CHECK (rating BETWEEN 1 AND 5) DEFAULT 5 NOT NULL,
        comment TEXT,
        created_at TIMESTAMP DEFAULT NOW () NOT NULL,
        updated_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

CREATE TABLE
    carts (
        id SERIAL PRIMARY KEY,
        user_id UUID CONSTRAINT fk_ca_user_id REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        product_id INT CONSTRAINT fk_ca_product_id REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        product_variant_id INT CONSTRAINT fk_ca_product_variant_id REFERENCES product_variants (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        CONSTRAINT uq_ca_product_product_variant UNIQUE (user_id, product_id, product_variant_id),
        quantity SMALLINT DEFAULT 1 NOT NULL,
        created_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

CREATE TABLE
    category_payments (id SERIAL PRIMARY KEY, name VARCHAR(50) NOT NULL);

-- 1. instant
-- 2. e-money
-- 3. transfer
CREATE TABLE
    payments (
        id SERIAL PRIMARY KEY,
        category_payment_id INT CONSTRAINT fk_py_category_payment_id REFERENCES category_payments (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        name VARCHAR(50) NOT NULL,
        img TEXT NOT NULL
    );

-- instant --
-- 1. Qris
-- e-money --
-- 2. GoPay
-- 3. DANA
-- transfer --
-- 5. Mandiri
-- 8. BCA
-- 9. BNI
-- 10. Permata
CREATE TABLE
    order_statuses (id SERIAL PRIMARY KEY, name VARCHAR(30) NOT NULL);

-- 1. menunggu pembayaran
-- 2. pembayaran dibatalkan
-- 3. batas waktu pembayaran habis
-- 4. menunggu konfirmasi seller
-- 5. dibatalkan seller
-- 6. dibatalkan pengguna
-- 7. produk sedang disiapkan
-- 8. produk sudah dikirim
-- 9. produk tidak diterima
-- 10. pesanan selesai
CREATE TABLE
    checkouts (
        id SERIAL PRIMARY KEY,
        user_id UUID CONSTRAINT fk_ch_user_id REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        product_id INT CONSTRAINT fk_ch_product_id REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        product_variant_id INT CONSTRAINT fk_ch_product_variant_id REFERENCES product_variants (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        CONSTRAINT uq_ch_order_product_product_variant UNIQUE (id, product_id, product_variant_id),
        quantity SMALLINT DEFAULT 1 NOT NULL,
        unit_price BIGINT NOT NULL,
        total_price BIGINT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

CREATE TABLE
    orders (
        id SERIAL PRIMARY KEY,
        user_id UUID CONSTRAINT fk_or_user_id REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        payment_id INT CONSTRAINT fk_or_payment_id REFERENCES payments (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        order_status_id INT CONSTRAINT fk_or_order_status_id REFERENCES order_statuses (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        product_id INT CONSTRAINT fk_or_product_id REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        product_variant_id INT CONSTRAINT fk_or_product_variant_id REFERENCES product_variants (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        CONSTRAINT uq_or_order_product_product_variant UNIQUE (id, product_id, product_variant_id),
        note VARCHAR(100),
        quantity SMALLINT DEFAULT 1 NOT NULL,
        unit_price BIGINT NOT NULL,
        total_price BIGINT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

ALTER TABLE orders
ADD COLUMN order_pretty_id VARCHAR(50) CONSTRAINT uq_or_order_pretty_id UNIQUE;

CREATE SEQUENCE order_pretty_id_seq START 1;

CREATE TABLE
    reminder_schedules (
        id UUID PRIMARY KEY,
        user_id UUID CONSTRAINT fk_rs_user_id REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        product_id INT CONSTRAINT fk_rs_product_id REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        product_variant_id INT CONSTRAINT fk_rs_product_variant_id REFERENCES product_variants (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        CONSTRAINT uq_rs_user_product_product_variant UNIQUE (user_id, product_id, product_variant_id),
        next_send DATE NOT NULL,
        next_warning_send DATE NOT NULL,
        next_remove DATE NOT NULL,
        last_sent_at DATE DEFAULT NOW () NOT NULL,
        is_over BOOLEAN DEFAULT FALSE NOT NULL,
        created_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

CREATE TABLE
    conversations (
        id UUID PRIMARY KEY,
        user_id UUID CONSTRAINT fk_co_user_id REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        seller_id UUID CONSTRAINT fk_co_seller_id REFERENCES sellers (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        CONSTRAINT uq_co_user_seller UNIQUE (user_id, seller_id),
        last_message_at TIMESTAMP DEFAULT NOW () NOT NULL,
        last_message_content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

CREATE TABLE
    messages (
        id SERIAL PRIMARY KEY,
        conversation_id UUID CONSTRAINT fk_m_conversation_id REFERENCES conversations (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        sender_id UUID NOT NULL, -- UUID of user or seller, no FK by design
        is_user BOOLEAN DEFAULT TRUE NOT NULL,
        content TEXT NOT NULL,
        sent_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

CREATE TABLE
    reports (
        id UUID PRIMARY KEY,
        user_id UUID CONSTRAINT fk_re_user_id REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        seller_id UUID CONSTRAINT fk_re_seller_id REFERENCES sellers (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        product_id INT CONSTRAINT fk_re_product_id REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        product_variant_id INT CONSTRAINT fk_re_product_variant_id REFERENCES product_variants (id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        description TEXT NOT NULL,
        target_role VARCHAR(10) CHECK (target_role IN ('admin', 'seller')),
        created_at TIMESTAMP DEFAULT NOW () NOT NULL
    );

-- TO BACKUP
-- BACKUP DATABASE testDB
-- TO DISK = '/path/to/backup/folder/file.bak';