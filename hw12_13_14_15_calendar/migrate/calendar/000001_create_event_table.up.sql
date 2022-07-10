-- CREATE SEQUENCE: public.event_id_seq

CREATE SEQUENCE IF NOT EXISTS public.event_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;


-- Table: public.event

CREATE TABLE IF NOT EXISTS public.event
(
    id bigint NOT NULL DEFAULT nextval('event_id_seq'::regclass),
    "title" text COLLATE pg_catalog."default" NOT NULL,
    user_id bigint NOT NULL,
    start_dt timestamp without time zone NOT NULL,
    end_dt timestamp without time zone NOT NULL,
    notif_dt timestamp without time zone,
    "desc" text COLLATE pg_catalog."default",
    CONSTRAINT event_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.event
    OWNER to postgres;

COMMENT ON TABLE public.event
    IS 'Событие - основная сущность';

COMMENT ON COLUMN public.event.title
    IS 'Заголовок';

COMMENT ON COLUMN public.event.user_id
    IS 'ID пользователя, владельца события';

COMMENT ON COLUMN public.event.start_dt
    IS 'Дата и время события';

COMMENT ON COLUMN public.event.end_dt
    IS 'дата и время окончания';

COMMENT ON COLUMN public.event.notif_dt
    IS 'Дата и время посылки уведомление';

COMMENT ON COLUMN public.event."desc"
    IS 'Описание события';


-- ALTER SEQUENCE: public.event_id_seq

ALTER SEQUENCE public.event_id_seq
    OWNER TO postgres;

ALTER SEQUENCE public.event_id_seq
    OWNED BY event.id;
