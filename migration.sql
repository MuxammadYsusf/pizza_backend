CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS pizza (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    cost FLOAT NOT NULL,
    type_id INT NOT NULL,
    FOREIGN KEY (type_id) REFERENCES types(id)
);

CREATE TABLE IF NOT EXISTS types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    date TIMESTAMP NOT NULL,
    is_ordered BOOLEAN NOT NULL,
    user_id INT NOT NULL,
    status TEXT NOT NULL,
    cart_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS order_item (
    id SERIAL PRIMARY KEY,
    pizza_id INT NOT NULL,
    cost FLOAT NOT NULL,
    quantity INT NOT NULL,
    order_id INT NOT NULL,
    FOREIGN KEY (pizza_id) REFERENCES pizza(id),
    FOREIGN KEY (order_id) REFERENCES orders(id)

);

CREATE TABLE IF NOT EXISTS cart (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    is_active BOOLEAN NOT NULL,
    total_cost FLOAT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS cart_item (
    id SERIAL PRIMARY KEY,
    pizza_id INT NOT NULL,
    pizza_type_id INT NOT NULL,
    cost INT NOT NULL,
    cart_id INT NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (cart_id) REFERENCES cart(id),
    FOREIGN KEY (pizza_id) REFERENCES pizza(id),
    FOREIGN KEY (pizza_type_id) REFERENCES types(id)
);
