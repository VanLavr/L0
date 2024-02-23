create database internship;

create table if not exists orders (
    order_uid text primary key,
    track_number text,
    entr text,

    -- delivery object in order
    delivery_id int,
    foreign key (delivery_id) references delivery(delivery_id),
    -- payment object in order
    t_action text,
    foreign key (t_action) references payment(t_action),

    locale text,
    internal_signature text,
    customer_id text,
    delivery_service text,
    shardkey text,
    sm_id int,
    date_created text,
    oof_shard text
);

create table if not exists delivery (
    delivery_id serial primary key,
    name text,
    phone text,
    zip text,
    city text,
    address text,
    region text,
    email text
);

create table if not exists payment (
    t_action text primary key,
    request_id text,
    currency text,
    provider text,
    amount float,
    payment_dt int,
    bank text,
    delivery_cost float,
    goods_total float,
    custom_fee float
);

create table if not exists items (
    chrt_id serial primary key,
    track_number text,
    price float,
    rid text,
    name text,
    sale float,
    size text,
    total_price float,
    nm_id int,
    brand text,
    status int
);

-- stands for many to many relation with orders and items
-- (many orders may include many items)
create table if not exists items_to_orders (
    order_uid text,
    chrt_id int,
    foreign key (order_uid) references orders(order_uid),
    foreign key (chrt_id) references items(chrt_id)
);