CREATE TABLE IF NOT EXISTS public.measurements
(
    id            serial PRIMARY KEY,
    room_id       integer NOT NULL REFERENCES public.rooms(id),
    device_id     integer NOT NULL REFERENCES public.devices(id),
    value         numeric(10, 2) NOT NULL,
    created_date  timestamptz NOT NULL,
    updated_date  timestamptz NOT NULL,
    deleted_date  timestamptz
);

-- Критично важливий індекс для адміністратора (ТЗ: перегляд за пристроєм за день/тиждень/місяць)
-- Він миттєво знаходитиме записи за умови: device_id = X AND deleted_date IS NULL AND created_date >= Y
CREATE INDEX IF NOT EXISTS idx_measurements_device_created 
ON public.measurements (device_id, created_date) 
WHERE deleted_date IS NULL;

-- Додатковий індекс для методу FindList (якщо вимірювання часто шукають по кімнатах)
CREATE INDEX IF NOT EXISTS idx_measurements_room 
ON public.measurements (room_id) 
WHERE deleted_date IS NULL;
