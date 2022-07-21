CREATE TABLE categories(
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        category  TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
CREATE TABLE products(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    category UUID NOT NULL REFERENCES categories(id),
    price FLOAT NOT NULL,
    feature TEXT NOT NULL,
    about TEXT NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
CREATE INDEX idx_products
    ON products(id,name,category,price,about,quantity);