CREATE TABLE in_the_box(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pdt_id UUID REFERENCES products(id),
    name TEXT NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
