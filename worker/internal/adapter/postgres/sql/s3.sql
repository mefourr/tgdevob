-- authT
DROP TABLE IF EXISTS public.s3_buckets;
DROP TABLE IF EXISTS public.s3_objects;

CREATE TABLE IF NOT EXISTS public.s3_buckets
(
    id      SERIAL PRIMARY KEY,
--     some bucket metadata
    user_id UUID NOT NULL,
    CONSTRAINT user_fk FOREIGN KEY (user_id) REFERENCES public.tg_users(id)
);
CREATE TABLE IF NOT EXISTS public.s3_objects
(
    id        SERIAL PRIMARY KEY,
    --     some bucket object metadata
    bucket_id INTEGER NOT NULL,
    CONSTRAINT bucket_fk FOREIGN KEY (bucket_id) REFERENCES public.s3_buckets(id)
);

-- ...
DROP TABLE IF EXISTS public.rec_msg_info;

CREATE TABLE IF NOT EXISTS public.rec_msg_info
(
    id        SERIAL PRIMARY KEY,
    --     some recognized message info
    object_id INTEGER NOT NULL,
    CONSTRAINT object_fk FOREIGN KEY (object_id) REFERENCES public.s3_objects(id)
);
