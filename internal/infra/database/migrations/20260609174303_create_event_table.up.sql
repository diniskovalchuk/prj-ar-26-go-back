CREATE TABLE IF NOT EXISTS public.events
(
    id              serial PRIMARY KEY,
    device_id       integer NOT NULL REFERENCES public.devices(id),
    room_id         integer NOT NULL REFERENCES public.rooms(id),
    action          varchar(10) NOT NULL, -- "on" або "off" (або "start"/"stop")
    created_date    timestamptz NOT NULL,
    updated_date    timestamptz NOT NULL,
    deleted_date    timestamptz
);

-- Індекс для аналітики в межах усього підприємства та окремих кімнат за часом
CREATE INDEX IF NOT EXISTS idx_events_room_created
ON public.events (room_id, created_date)
WHERE deleted_date IS NULL;

CREATE INDEX IF NOT EXISTS idx_events_device_created
ON public.events (device_id, created_date)
WHERE deleted_date IS NULL;
