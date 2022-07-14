CREATE TABLE cart_products(
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        product_id UUID references products(id),
        cart_id UUID references carts(id),
        quantity INT not null,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);