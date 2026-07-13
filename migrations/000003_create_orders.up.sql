CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    customer_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    total NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_customer
        FOREIGN KEY(customer_id) 
        REFERENCES customers(id)
        ON DELETE CASCADE
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    order_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    price NUMERIC(10, 2) NOT NULL,
    CONSTRAINT fk_order
        FOREIGN KEY(order_id) 
        REFERENCES orders(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_product
        FOREIGN KEY(product_id) 
        REFERENCES products(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
