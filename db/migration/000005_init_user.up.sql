DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM public.users WHERE id = 1) THEN
        INSERT INTO public.users (id, code, username, first_name, last_name, email, phone, description, hashed_password)
        VALUES (1, 'U000001', 'administrator', 'Administrator', 'System', 'admin@localhost.com', '+66111111111', NULL, '$2a$10$v/d2ZXL5z92Y47Y.y5Y.Q.WYiwb/wWZvc2SPpxcuBlKKbPRAe2yD.');
    END IF;
END $$;
