INSERT INTO public.user_role (user_id, role_id, created_at) 
SELECT 1, id, NOW()::TIMESTAMP FROM public.role;