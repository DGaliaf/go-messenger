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

END;