CREATE TABLE orders(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID references products(id),
    user_id UUID references users(id),
    address_id UUID references addresses(id),
    cost FLOAT NOT NULL,
    payment_by_card UUID DEFAULT NULL references cards(id),
    payment_by_cod bool DEFAULT FALSE,
    date_of_delivery TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
