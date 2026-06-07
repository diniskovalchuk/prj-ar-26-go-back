CREATE TABLE IF NOT EXISTS public.devices
(
    Id              serial PRIMARY KEY,
    OrganizationId  integer NOT NULL REFERENCES public.organizations(id),
    RoomId          varchar(250) NOT NULL,
    GUID            uuid NOT NULL,
    InventoryNumber varchar(250),
    SerialNumber    varchar(250),
    Characteristics text,
    Category        varchar(250),
    Units           varchar(250),
    PowerConsumption numeric(10, 2),
    CreatedDate     timestamptz NOT NULL,
    UpdatedDate     timestamptz NOT NULL,
    DeletedDate     timestamptz
);