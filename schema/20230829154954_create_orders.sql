-- +goose Up
-- +goose StatementBegin
CREATE TABLE delivery_info
(
    id SERIAL NOT NULL UNIQUE ,
    name      VARCHAR(255) not null ,
    phone     VARCHAR(20) not null ,
    zip       VARCHAR(20) not null ,
    city      VARCHAR(255),
    address   TEXT not null ,
    region    VARCHAR(255) not null ,
    email     VARCHAR(255) not null unique
);
CREATE TABLE payment_info
(
    id SERIAL NOT NULL UNIQUE ,
    transaction   VARCHAR(255) NOT NULL UNIQUE ,
    request_id    VARCHAR(255),
    currency      VARCHAR(10) NOT NULL ,
    provider      VARCHAR(255) NOT NULL ,
    amount        INT NOT NULL ,
    payment_dt    INT NOT NULL ,
    bank          VARCHAR(255) NOT NULL ,
    delivery_cost INT,
    goods_total   INT,
    custom_fee    INT
);
CREATE TABLE items
(
    id serial not null unique ,
    chrt_id      INT NOT NULL UNIQUE ,
    track_number VARCHAR(255) NOT NULL UNIQUE ,
    price        INT NOT NULL ,
    rid          VARCHAR(255) NOT NULL UNIQUE ,
    name         VARCHAR(255) NOT NULL ,
    sale         INT,
    size         VARCHAR(10),
    total_price  INT,
    nm_id        INT NOT NULL UNIQUE ,
    brand        VARCHAR(255) NOT NULL ,
    status       INT
);
CREATE TABLE orders
(
    order_uid        VARCHAR(255) NOT NULL unique ,
    track_number     VARCHAR(255)NOT NULL unique ,
    entry            VARCHAR(255) NOT NULL ,
    delivery_id INT NOT NULL UNIQUE REFERENCES delivery_info(id) on delete cascade,
    payment_id INT NOT NULL UNIQUE REFERENCES payment_info(id) on delete cascade,
    locale           VARCHAR(10) NOT NULL ,
    internal_signature varchar(255),
    customer_id      VARCHAR(255) NOT NULL UNIQUE ,
    delivery_service VARCHAR(255) NOT NULL ,
    shard_key         VARCHAR(10) NOT NULL ,
    sm_id            INT NOT NULL ,
    date_created     TIMESTAMP default now(),
    oof_shard        VARCHAR(10)
);

CREATE TABLE order_items
(
    order_id VARCHAR(255) references  orders(order_uid) on delete cascade,
    item_id  INT REFERENCES items(id) on delete cascade
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE order_items;
DROP TABLE orders;
DROP TABLE items;
DROP TABLE payment_info;
DROP TABLE delivery_info;



-- +goose StatementEnd
