CREATE TABLE IF NOT EXISTS likes(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
) 