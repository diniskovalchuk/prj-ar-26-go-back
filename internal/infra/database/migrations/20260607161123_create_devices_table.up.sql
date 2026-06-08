CREATE TABLE IF NOT EXISTS public.devices
(
    id              serial PRIMARY KEY,
    organization_id  integer NOT NULL REFERENCES public.organizations(id),
    room_id          integer NOT NULL REFERENCES public.rooms(id),
    guid            uuid NOT NULL DEFAULT gen_random_uuid(),
    inventory_number varchar(250),
    serial_number    varchar(250),
    characteristics text,
    category        varchar(250),
    units           varchar(250),
    power_consumption numeric(10, 2),
    created_date    timestamptz NOT NULL,
    updated_date    timestamptz NOT NULL,
    deleted_date    timestamptz
);