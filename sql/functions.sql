CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Yes, the uuid_generate_v4() function will generate a random and unique UUID (Universally Unique Identifier) 
-- each time it is called. This function is designed to produce a new, 
-- randomly generated UUID that is highly likely to be unique across all systems and all time.
-- Therefore, using uuid_generate_v4() within the generate_unique_id() function 
-- will indeed produce a random unique string ID each time the function is called.

-- generating unique string ids based on pseudo-random values.
CREATE OR REPLACE FUNCTION generate_unique_id()
RETURNS text AS $$
BEGIN
    RETURN lower(translate(uuid_generate_v4()::text, '-', ''));
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_order_by_id(order_id_to_find text)
RETURNS TABLE (
    order_uid text, 
    track_number text, 
    entr text, 
    delivery_name text, 
    delivery_phone text, 
    delivery_zip text, 
    delivery_city text, 
    delivery_address text, 
    delivery_region text, 
    delivery_email text, 
    payment_transaction text, 
    payment_request_id text, 
    payment_currency text, 
    payment_provider text, 
    payment_amount float, 
    payment_date int, 
    payment_bank text, 
    delivery_cost float, 
    goods_total float, 
    custom_fee float, 
    item_chrt_id int, 
    item_track_number text, -- added
    item_price float, 
    item_rid text, 
    item_name text, 
    item_sale float, 
    item_size text, 
    item_total_price float, 
    item_nm_id int, 
    item_brand text, 
    item_status int, 
    locale text, 
    internal_signature text, 
    customer_id text, 
    delivery_service text, 
    shardkey text, 
    sm_id int, 
    date_created text, 
    oof_shard text
) AS $$
BEGIN
    RETURN QUERY 
    SELECT 
        o.order_uid, 
        o.track_number, 
        o.entr AS entry, 
        d.name AS delivery_name, 
        d.phone AS delivery_phone, 
        d.zip AS delivery_zip, 
        d.city AS delivery_city, 
        d.address AS delivery_address, 
        d.region AS delivery_region, 
        d.email AS delivery_email, 
        p.t_action AS payment_transaction, 
        p.request_id AS payment_request_id, 
        p.currency AS payment_currency, 
        p.provider AS payment_provider, 
        p.amount AS payment_amount, 
        p.payment_dt AS payment_date, 
        p.bank AS payment_bank, 
        p.delivery_cost, 
        p.goods_total, 
        p.custom_fee, 
        i.chrt_id AS item_chrt_id, 
        i.track_number AS item_track_number, -- added
        i.price AS item_price, 
        i.rid AS item_rid, 
        i.name AS item_name, 
        i.sale AS item_sale, 
        i.size AS item_size, 
        i.total_price AS item_total_price, 
        i.nm_id AS item_nm_id, 
        i.brand AS item_brand, 
        i.status AS item_status, 
        o.locale, 
        o.internal_signature, 
        o.customer_id, 
        o.delivery_service, 
        o.shardkey, 
        o.sm_id, 
        o.date_created, 
        o.oof_shard
    FROM orders o
    INNER JOIN delivery d ON o.delivery_id = d.delivery_id
    INNER JOIN payment p ON o.t_action = p.t_action
    INNER JOIN items i ON o.track_number = i.track_number
    WHERE o.order_uid = order_id_to_find;
END;
$$ LANGUAGE plpgsql;