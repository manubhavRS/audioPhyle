CREATE TABLE addresses(
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID references users(id),
        address TEXT NOT NULL,
        landmark TEXT,
        phone_number TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);