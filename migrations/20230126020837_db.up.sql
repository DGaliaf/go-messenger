BEGIN;

CREATE TABLE IF NOT EXISTS public.user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS public.chat (
    id SERIAL PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS public.user_chat (
    user_id UUID REFERENCES public.user(id),
    chat_id INT REFERENCES public.chat(id)
);

CREATE TABLE IF NOT EXISTS public.message (
    id SERIAL PRIMARY KEY NOT NULL,
    chat_id INT REFERENCES public.chat(id),
    author_id UUID REFERENCES public.user(id),
    text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

-- INSERT INTO public.user(username) VALUES ('test_user_2');
-- INSERT INTO public.user(username) VALUES ('test_user_3');
-- INSERT INTO public.chat(name) VALUES ('test_chat_2');
-- INSERT INTO public.user_chat(user_id, chat_id) VALUES ('d5405cfa-5110-4276-bbbf-4f33dfd88ed0', 2);
-- INSERT INTO public.user_chat(user_id, chat_id) VALUES ('23f6e1d9-4e99-4297-b811-641a9210a9e0', 2);

-- SELECT * FROM (
--     (public.user_chat INNER JOIN public.chat ON public.chat.id = public.user_chat.chat_id)
-- );

--     INNER JOIN public.user ON public.user.id = public.user_chat.user_id);

-- SELECT public.user.id, public.user.username, public.user.created_at,
--        public.chat.id, public.chat.name, public.chat.created_at
-- FROM public.user
-- JOIN public.user_chat ON public.user.id = public.user_chat.user_id
-- JOIN public.chat ON public.chat.id = public.user_chat.chat_id;

END;