CREATE OR REPLACE FUNCTION check_user_blocked(current_user_id UUID, parent_id UUID)
RETURNS VOID AS $$
DECLARE
    post_user_id UUID;
BEGIN
    SELECT user_id INTO post_user_id FROM posts WHERE id = parent_id;
    
    IF EXISTS (
        SELECT 1
        FROM blocking
        WHERE (user_id = current_user_id AND blocking_id = post_user_id) 
           OR (user_id = post_user_id AND blocking_id = current_user_id)
    ) THEN
        RAISE EXCEPTION 'User cannot access this post';
    END IF;
END;
$$ LANGUAGE plpgsql;
