DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM public.users WHERE id = 1) THEN
        DELETE FROM public.users WHERE id = 1;
    END IF;
END $$;
