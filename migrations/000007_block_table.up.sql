CREATE TABLE IF NOT EXISTS blocking(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
    user_id UUID NOT NULL,
    blocking_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ,
    FOREIGN KEY (blocking_id) REFERENCES users(id) ON DELETE CASCADE
) 


