-- Users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at BIGINT NOT NULL DEFAULT (extract(epoch FROM now()) * 1000),
    UNIQUE(email),
    CHECK (role IN ('ADMIN','USER'))
);

CREATE TABLE IF NOT EXISTS types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at BIGINT NOT NULL DEFAULT (extract(epoch FROM now()) * 1000),
    UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS pizza (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    cost INT NOT NULL,
    photo TEXT NOT NULL,
    type_id INT NOT NULL REFERENCES types(id) ON DELETE CASCADE,
    created_at BIGINT NOT NULL DEFAULT (extract(epoch FROM now()) * 1000),
    CONSTRAINT uq_pizza_type_name UNIQUE (type_id, name)
);

CREATE TABLE IF NOT EXISTS cart (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at BIGINT NOT NULL DEFAULT (extract(epoch FROM now()) * 1000)
);

CREATE TABLE IF NOT EXISTS cart_item (
    id SERIAL PRIMARY KEY,
    pizza_id INT NOT NULL REFERENCES pizza(id) ON DELETE CASCADE,
    pizza_type_id INT NOT NULL REFERENCES types(id) ON DELETE CASCADE,
    cost INT NOT NULL,
    cart_id INT NOT NULL REFERENCES cart(id) ON DELETE CASCADE,
    quantity INT NOT NULL,
    created_at BIGINT NOT NULL DEFAULT (extract(epoch FROM now()) * 1000),
    CONSTRAINT ck_cart_item_quantity CHECK (quantity > 0),
    CONSTRAINT uq_cart_item_line UNIQUE (cart_id, pizza_id, pizza_type_id)
);

CREATE TABLE IF NOT EXISTS order (
    id SERIAL PRIMARY KEY,
    order_time BIGINT NOT NULL,
    status TEXT NOT NULL,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cart_id INT NOT NULL REFERENCES cart(id),
    created_at BIGINT NOT NULL DEFAULT (extract(epoch FROM now()) * 1000),
    CONSTRAINT ck_order_status CHECK (status IN ('NEW','PAID','IN_PROGRESS','DELIVERING','CANCELED','DONE'))
);

CREATE TABLE IF NOT EXISTS order_item (
    id SERIAL PRIMARY KEY,
    pizza_id INT NOT NULL REFERENCES pizza(id) ON DELETE CASCADE,
    cost INT NOT NULL,
    quantity INT NOT NULL,
    order_id INT NOT NULL REFERENCES order(id) ON DELETE CASCADE,
    created_at BIGINT NOT NULL DEFAULT (extract(epoch FROM now()) * 1000),
    CONSTRAINT ck_order_item_quantity CHECK (quantity > 0),
    CONSTRAINT uq_order_item_line UNIQUE (order_id, pizza_id)
);

CREATE TABLE IF NOT EXISTS session (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL,
    expired_time BIGINT NOT NULL,
    created_at BIGINT NOT NULL DEFAULT (extract(epoch FROM now()) * 1000),
    CHECK (role IN ('ADMIN','USER'))
);