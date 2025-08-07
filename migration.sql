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

CREATE OR REPLACE FUNCTION reduce_total_cost()
RETURNS TRIGGER AS $$
DECLARE
    _cost FLOAT;
BEGIN
    SELECT cost INTO _cost FROM cart WHERE id = NEW.id;
    NEW.total_cost := OLD.total_cost - _cost;
    RETURN new;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER reduce_total_cost
BEFORE UPDATE ON cart
FOR EACH ROW EXECUTE PROCEDURE reduce_total_cost();



CREATE OR REPLACE FUNCTION cancel_from_cart()
RETURNS TRIGGER AS $$
BEGIN
    NEW.cost := OLD.cost - (OLD.cost - NEW.cost);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cencel_from_cart
BEFORE UPDATE ON cart_item
FOR EACH ROW EXECUTE PROCEDURE cancel_from_cart();
