-- Table: public.event

CREATE TABLE IF NOT EXISTS public.event
(
    id bigint NOT NULL DEFAULT nextval('event_id_seq'::regclass),
    title character varying(100)[] COLLATE pg_catalog."default" NOT NULL,
    start_time time without time zone NOT NULL,
    end_time time without time zone NOT NULL,
    describtion text COLLATE pg_catalog."default",
    user_id bigint NOT NULL,
    notification_time timestamp without time zone,
    CONSTRAINT event_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.event
    OWNER to postgres;

COMMENT ON TABLE public.event
    IS 'Событие - основная сущность';

COMMENT ON COLUMN public.event.id
    IS 'уникальный идентификатор события';

COMMENT ON COLUMN public.event.title
    IS 'Заголовок';

COMMENT ON COLUMN public.event.start_time
    IS 'Дата и время события';

COMMENT ON COLUMN public.event.end_time
    IS 'дата и время окончания';

COMMENT ON COLUMN public.event.describtion
    IS 'Описание события';

COMMENT ON COLUMN public.event.user_id
    IS 'ID пользователя, владельца события';

COMMENT ON COLUMN public.notification_time
    IS 'За сколько времени высылать уведомление';