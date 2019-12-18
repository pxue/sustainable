CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name text NOT NULL DEFAULT ''::text,
    factory_link text NOT NULL DEFAULT '':text
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX ON brands(id int4_ops);

-- Table Definition ----------------------------------------------
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    code text NOT NULL DEFAULT ''::text,
    category text NOT NULL,
    name text NOT NULL DEFAULT ''::text,
    materials jsonb NOT NULL DEFAULT '{}'::jsonb,
    price numeric(7,2) NOT NULL DEFAULT 0
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX ON products(id int4_ops);
CREATE UNIQUE INDEX ON products(code text_ops);

-- Table Definition ----------------------------------------------

CREATE TABLE suppliers (
    id SERIAL PRIMARY KEY,
    name text NOT NULL DEFAULT ''::text,
    country text NOT NULL DEFAULT ''::text,
    hm_id text NOT NULL DEFAULT ''::text
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX ON suppliers(id int4_ops);
CREATE UNIQUE INDEX ON suppliers(hm_id text_ops);

-- Table Definition ----------------------------------------------

CREATE TABLE factories (
    id SERIAL PRIMARY KEY,
    supplier_id bigint REFERENCES suppliers(id) ON DELETE CASCADE,
    name text NOT NULL DEFAULT ''::text,
    country text NOT NULL DEFAULT ''::text,
    address text NOT NULL DEFAULT ''::text,
    code text NOT NULL DEFAULT ''::text,
    lat numeric(10,7) NOT NULL DEFAULT '0'::numeric,
    lon numeric(10,7) NOT NULL DEFAULT '0'::numeric
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX ON factories(id int4_ops);
CREATE UNIQUE INDEX ON factories(hm_id text_ops);

-- Table Definition ----------------------------------------------

CREATE TABLE product_suppliers (
    product_id bigint REFERENCES products(id),
    supplier_id bigint REFERENCES suppliers(id),
    factory_id bigint REFERENCES factories(id) ON DELETE CASCADE
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX ON product_suppliers(product_id int8_ops,supplier_id int8_ops,factory_id int8_ops);

-- Table Definition ----------------------------------------------

CREATE TABLE materials (
    id SERIAL PRIMARY KEY,
    name text NOT NULL DEFAULT ''::text,
    "type" text NOT NULL DEFAULT ''::text,
    msi_score bigint NOT NULL DEFAULT '0'::bigint
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX ON materials(id int4_ops);


CREATE TABLE bci_members (
    id SERIAL PRIMARY KEY,
    name text NOT NULL DEFAULT ''::text,
    category text NOT NULL DEFAULT ''::text,
    country text NOT NULL DEFAULT ''::text,
    website text NOT NULL DEFAULT ''::text,
    since timestamptz NOT NULL DEFAULT now()
);

-- Indices -------------------------------------------------------
