CREATE TABLE product_assets(
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        pdt_id UUID REFERENCES products(id),
        name TEXT UNIQUE NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
