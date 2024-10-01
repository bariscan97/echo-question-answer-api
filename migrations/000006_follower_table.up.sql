CREATE TABLE IF NOT EXISTS following(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
    user_id UUID NOT NULL,
    following_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ,
    FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE
) 