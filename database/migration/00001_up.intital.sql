
begin;

create type asset_type as enum (
    'laptop',
    'keyboard',
    'mouse',
    'mobile',
    'hardware'
    );

create type asset_status as enum (
    'available',
    'assigned',
    'in_service',
    'for_repair',
    'damaged'
    );

create type user_role as enum (
    'admin',
    'employee',
    'intern',
    'freelancer',
    'manager'
    );

create type owner_type as enum (
    'client',
    'company'
    );

create table employee
(
    id            uuid primary key default gen_random_uuid(),
    name          varchar(100) not null,
    email         text  not null  ,
    role          user_role default'employee',
    phone_no      text  not null unique ,
    password_hash text  not null,
    created_at    timestamp  default now(),
    archived_at   timestamptz
);
create table user_sessions
(
    id          uuid primary key default gen_random_uuid(),
    emp_id      uuid references employee(id),
    created_at  timestamptz default now(),
    archived_at timestamptz
);


create unique index idx_unique_email on employee (email) where archived_at is null;

create table assets
(
    id             uuid primary key default gen_random_uuid(),
    emp_id         uuid references employee(id),
    -- name           varchar(100) not null,
    brand          text         not null,
    model          varchar(100) not null,
    serial_no      text         not null,
    type           asset_type   not null,
    status         asset_status     default 'available',
    purchased_at timestamptz not null ,
    warranty_start timestamptz,
    warranty_end   timestamptz,
    owner          owner_type   not null,
    archived_at    timestamptz,
    Note           text,
    constraint chk_warranty check (
        (warranty_start is null and warranty_end is null)
            or
        (warranty_start is not null and warranty_end is not null and warranty_end > warranty_start)
        )
);


create unique index idx_unique_serial on assets (serial_no) where archived_at is null;
create type history as enum(
    'assigned', 'available','repair'
    );
create table asset_history
(
    id            uuid primary key default gen_random_uuid(),
    type history not null,
    asset_id      uuid references assets(id),
    assigned_to   uuid references employee(id),
    assigned_on   timestamptz not null,
    returned_on   timestamptz,
    return_status text
);

create table laptop
(
    asset_id  uuid primary key references assets(id),
    processor text,
    ram       int,
    storage   int,
    os        varchar(50),
    charger   varchar(50)
);

create table keyboard
(
    asset_id uuid primary key references assets(id),
    layout   text
);

create table mouse
(
    asset_id    uuid primary key references assets(id),
    dpi         int,
    is_wireless boolean default false
);

create table mobile
(
    asset_id uuid primary key references assets(id),
    os       text,
    ram      int,
    storage  int,
    charger  varchar(50)
);
create table if not exists hardware (
                                        asset_id uuid primary key references assets(id),
                                        storage  int not null
);


-- create or replace function enforce_laptop_type()
--     returns trigger as $$
-- begin
--     if not exists (
--         select 1 from assets
--         where id = new.asset_id
--           and type = 'laptop'
--     ) then
--         raise exception 'Asset type mismatch: not a laptop';
--     end if;
--
--     return new;
-- end;
-- $$ language plpgsql;
--
-- create trigger trg_laptop_type_check
--     before insert or update on laptop
--     for each row
-- execute function enforce_laptop_type();


create or replace function enforce_asset_type()
    returns trigger as $$
declare
    expected asset_type := TG_ARGV[0]::asset_type;
begin
    if not exists (
        select 1 from assets where id = new.asset_id and type = expected
    ) then
        raise exception 'Asset type mismatch: expected %', expected;
    end if;
    return new;
end;
$$ language plpgsql;

create trigger trg_mouse_type_check before insert or update on mouse
    for each row execute function enforce_asset_type('mouse');

create trigger trg_keyboard_type_check before insert or update on keyboard
    for each row execute function enforce_asset_type('keyboard');

commit;