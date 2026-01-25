-- 1
CREATE TABLE public.users (
	id serial4 NOT NULL,
	email varchar(255) NULL,
	"password" varchar(255) NULL,
	"role" varchar(50) NULL,
	tenant varchar(100) NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);
-- 2
CREATE TABLE public.products (
	id serial4 NOT NULL,
	code varchar(100) NULL,
	"name" varchar(255) NULL,
	het_price numeric NULL,
	cost_price numeric NULL,
	avg_price numeric NULL,
	stock_on_hand int4 NULL,
	stock_allocation int4 NULL,
	stock_back_order int4 NULL,
	stock_packing int4 NULL,
	stock_on_purchase int4 NULL,
	stock_on_receive int4 NULL,
	lead_time_days int4 NULL,
	tenant varchar(100) NULL,
	supplier_id int8 NULL,
	created_at timestamp NULL,
	CONSTRAINT products_pkey PRIMARY KEY (id)
);
-- 3
CREATE TABLE public.product_locations (
	id serial4 NOT NULL,
	product_id int8 NULL,
	location_id int8 NULL,
	location_code varchar(100) NULL,
	qtty int4 NULL,
	CONSTRAINT product_locations_pkey PRIMARY KEY (id)
);


-- public.product_locations foreign keys

ALTER TABLE public.product_locations ADD CONSTRAINT product_locations_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);
-- 4
CREATE TABLE public.demand (
	id serial4 NOT NULL,
	product_id int8 NULL,
	"month" varchar(20) NULL,
	n_qty int8 NULL,
	n1_qty int8 NULL,
	n2_qty int8 NULL,
	n3_qty int8 NULL,
	n4_qty int8 NULL,
	n5_qty int8 NULL,
	n6_qty int8 NULL,
	n7_qty int8 NULL,
	n8_qty int8 NULL,
	n9_qty int8 NULL,
	n10_qty int8 NULL,
	n11_qty int8 NULL,
	n12_qty int8 NULL,
	CONSTRAINT demand_pkey PRIMARY KEY (id)
);
-- 5
CREATE TABLE public.locations (
	id serial4 NOT NULL,
	code varchar(100) NULL,
	"name" varchar(255) NULL,
	tenant varchar(100) NULL,
	CONSTRAINT locations_pkey PRIMARY KEY (id)
);
-- 6
CREATE TABLE public.configs (
	id serial4 NOT NULL,
	"name" varchar(255) NULL,
	description varchar(100) NULL,
	value varchar(255) NULL,
	tenant varchar(100) NULL,
	CONSTRAINT configs_pkey PRIMARY KEY (id)
);
-- 7
CREATE TABLE public.customers (
	id serial4 NOT NULL,
	"name" varchar(255) NULL,
	email varchar(255) NULL,
	phone varchar(50) NULL,
	address text NULL,
	discount_percent int8 NULL,
	term_of_payment varchar(100) NULL,
	cancel_on_back_order bool NULL,
	tenant varchar(100) NULL,
	CONSTRAINT customers_pkey PRIMARY KEY (id)
);
-- 8
CREATE TABLE public.suppliers (
	id serial4 NOT NULL,
	"name" varchar(255) NULL,
	email varchar(255) NULL,
	phone varchar(50) NULL,
	address text NULL,
	discount_percent int8 NULL,
	term_of_payment varchar(100) NULL,
	tenant varchar(100) NULL,
	lead_time_days int4 NULL,
	CONSTRAINT suppliers_pkey PRIMARY KEY (id)
);
-- 9
CREATE TABLE public.purchase_orders (
	id serial4 NOT NULL,
	supplier_id int8 NULL,
	order_date timestamp NULL,
	amount numeric NULL,
	discount int4 NULL,
	total numeric NULL,
	tenant varchar(100) NULL,
	purchase_number varchar(50) NULL,
	CONSTRAINT purchase_orders_pkey PRIMARY KEY (id)
);
-- 10
CREATE TABLE public.purchase_order_items (
	id serial4 NOT NULL,
	supplier_id int8 NULL,
	purchase_order_id int8 NULL,
	product_id int8 NULL,
	order_qty int4 NULL,
	receive_order_qty int4 NULL,
	stocked_order_qty int4 NULL,
	price numeric NULL,
	total numeric NULL,
	sub_total numeric NULL,
	discount numeric NULL,
	CONSTRAINT purchase_order_items_pkey PRIMARY KEY (id)
);


-- public.purchase_order_items foreign keys

ALTER TABLE public.purchase_order_items ADD CONSTRAINT purchase_order_items_purchase_order_id_fkey FOREIGN KEY (purchase_order_id) REFERENCES public.purchase_orders(id);
-- 11
CREATE TABLE public.receive_orders (
	id serial4 NOT NULL,
	supplier_id int8 NULL,
	receive_date timestamp NULL,
	status varchar(50) NULL,
	tenant varchar(100) NULL,
	due_date timestamp NULL,
	payment_status varchar(20) NULL,
	invoice_number varchar(20) NULL,
	payment_date timestamp NULL,
	amount numeric NULL,
	total_amount numeric NULL,
	CONSTRAINT receive_orders_pkey PRIMARY KEY (id)
);
-- 12
CREATE TABLE public.receive_order_items (
	id serial4 NOT NULL,
	receive_order_id int8 NULL,
	purchase_order_id int8 NULL,
	product_id int8 NULL,
	product_code varchar(100) NULL,
	product_name varchar(255) NULL,
	product_location varchar(100) NULL,
	receive_order_qty int4 NULL,
	purchase_price numeric NULL,
	price numeric NULL,
	CONSTRAINT receive_order_items_pkey PRIMARY KEY (id)
);


-- public.receive_order_items foreign keys

ALTER TABLE public.receive_order_items ADD CONSTRAINT receive_order_items_receive_order_id_fkey FOREIGN KEY (receive_order_id) REFERENCES public.receive_orders(id);

-- 13
CREATE TABLE public.account_payable (
	id bigserial NOT NULL,
	receive_id int8 NULL,
	amount numeric NULL,
	payment_date timestamp NULL,
	payment_method varchar(30) NULL,
	reference_number varchar(100) NULL,
	notes varchar(255) NULL,
	tenant varchar(100) NULL,
	invoice_number varchar(100) NULL
);
-- 14
CREATE TABLE public.sales_orders (
	id serial4 NOT NULL,
	customer_id int8 NULL,
	order_date timestamp NULL,
	amount numeric NULL,
	discount int4 NULL,
	total numeric NULL,
	tenant varchar(100) NULL,
	CONSTRAINT sales_orders_pkey PRIMARY KEY (id)
);
-- 15
CREATE TABLE public.sales_order_items (
	id serial4 NOT NULL,
	sales_order_id int8 NULL,
	product_id int8 NULL,
	order_qty int4 NULL,
	back_order_qty int4 NULL,
	allocation_order_qty int4 NULL,
	packing_order_qty int4 NULL,
	invoice_order_qty int4 NULL,
	price numeric NULL,
	total numeric NULL,
	customer_id int8 NULL,
	CONSTRAINT sales_order_items_pkey PRIMARY KEY (id)
);
-- 16
CREATE TABLE public.packing_orders (
	id serial4 NOT NULL,
	customer_id int8 NULL,
	packing_date timestamp NULL,
	status varchar(50) NULL,
	tenant varchar(100) NULL,
	CONSTRAINT packing_orders_pkey PRIMARY KEY (id)
);
-- 17
CREATE TABLE public.packing_order_items (
	id serial4 NOT NULL,
	packing_order_id int8 NULL,
	sales_order_id int8 NULL,
	product_id int8 NULL,
	product_code varchar(100) NULL,
	product_name varchar(255) NULL,
	product_location varchar(100) NULL,
	packing_order_qty int4 NULL,
	customer_id int8 NULL,
	CONSTRAINT packing_order_items_pkey PRIMARY KEY (id)
);


-- public.packing_order_items foreign keys

ALTER TABLE public.packing_order_items ADD CONSTRAINT packing_order_items_packing_order_id_fkey FOREIGN KEY (packing_order_id) REFERENCES public.packing_orders(id);
-- 18
CREATE TABLE public.invoice_orders (
	id serial4 NOT NULL,
	packing_order_id int8 NULL,
	customer_id int8 NULL,
	customer_name varchar(255) NULL,
	invoice_date timestamp NULL,
	amount numeric NULL,
	discount int4 NULL,
	total_amount numeric NULL,
	status varchar(50) NULL,
	tenant varchar(100) NULL,
	due_date timestamp NULL,
	payment_status varchar(20) NULL,
	invoice_number varchar(50) NULL,
	payment_date timestamp NULL,
	CONSTRAINT invoice_orders_pkey PRIMARY KEY (id)
);
-- 19
CREATE TABLE public.invoice_order_items (
	id serial4 NOT NULL,
	invoice_order_id int8 NULL,
	product_code varchar(100) NULL,
	product_name varchar(255) NULL,
	quantity int4 NULL,
	price numeric NULL,
	total numeric NULL,
	sales_order_id int8 NULL,
	CONSTRAINT invoice_order_items_pkey PRIMARY KEY (id)
);


-- public.invoice_order_items foreign keys

ALTER TABLE public.invoice_order_items ADD CONSTRAINT invoice_order_items_invoice_order_id_fkey FOREIGN KEY (invoice_order_id) REFERENCES public.invoice_orders(id);
-- 20
CREATE TABLE public.account_receivable (
	id bigserial NOT NULL,
	invoice_id int8 NULL,
	amount numeric NULL,
	payment_date timestamp NULL,
	payment_method varchar(30) NULL,
	reference_number varchar(100) NULL,
	notes varchar(255) NULL,
	tenant varchar(100) NULL,
	invoice_number varchar(100) NULL
);