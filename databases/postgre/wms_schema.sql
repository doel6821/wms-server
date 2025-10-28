
CREATE TABLE customers (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(50),
  address TEXT,
  discount_percent BIGINT,
  term_of_payment VARCHAR(100),
  cancel_on_back_order BOOLEAN,
  tenant VARCHAR(100)
);

CREATE TABLE demand (
  id SERIAL PRIMARY KEY,
  product_id BIGINT,
  month VARCHAR(20),
  n_qty BIGINT,
  n1_qty BIGINT,
  n2_qty BIGINT,
  n3_qty BIGINT,
  n4_qty BIGINT,
  n5_qty BIGINT,
  n6_qty BIGINT,
  n7_qty BIGINT,
  n8_qty BIGINT,
  n9_qty BIGINT,
  n10_qty BIGINT,
  n11_qty BIGINT,
  n12_qty BIGINT
);

CREATE TABLE invoice_orders (
  id SERIAL PRIMARY KEY,
  packing_order_id BIGINT,
  customer_id BIGINT,
  customer_name VARCHAR(255),
  invoice_date TIMESTAMP,
  amount NUMERIC,
  discount INTEGER,
  total_amount NUMERIC,
  status VARCHAR(50),
  tenant VARCHAR(100)
);

CREATE TABLE invoice_order_items (
  id SERIAL PRIMARY KEY,
  invoice_order_id BIGINT REFERENCES invoice_orders(id),
  product_code VARCHAR(100),
  product_name VARCHAR(255),
  quantity INTEGER,
  price NUMERIC,
  total NUMERIC
);

CREATE TABLE locations (
  id SERIAL PRIMARY KEY,
  code VARCHAR(100),
  name VARCHAR(255),
  tenant VARCHAR(100)
);

CREATE TABLE packing_orders (
  id SERIAL PRIMARY KEY,
  customer_id BIGINT,
  packing_date TIMESTAMP,
  status VARCHAR(50),
  tenant VARCHAR(100)
);

CREATE TABLE packing_order_items (
  id SERIAL PRIMARY KEY,
  packing_order_id BIGINT REFERENCES packing_orders(id),
  sales_order_id BIGINT,
  product_id BIGINT,
  product_code VARCHAR(100),
  product_name VARCHAR(255),
  product_location VARCHAR(100),
  packing_order_qty INTEGER
);

CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  code VARCHAR(100),
  name VARCHAR(255),
  het_price NUMERIC,
  cost_price NUMERIC,
  avg_price NUMERIC,
  stock_on_hand INTEGER,
  stock_allocation INTEGER,
  stock_back_order INTEGER,
  stock_packing INTEGER,
  stock_on_purchase INTEGER,
  stock_on_receive INTEGER,
  lead_time_days INTEGER,
  tenant VARCHAR(100)
);

CREATE TABLE product_locations (
  id SERIAL PRIMARY KEY,
  product_id BIGINT REFERENCES products(id),
  location_id BIGINT,
  location_code VARCHAR(100),
  qtty INTEGER
);

CREATE TABLE purchase_orders (
  id SERIAL PRIMARY KEY,
  supplier_id BIGINT,
  order_date TIMESTAMP,
  amount NUMERIC,
  discount INTEGER,
  total NUMERIC,
  tenant VARCHAR(100)
);

CREATE TABLE purchase_order_items (
  id SERIAL PRIMARY KEY,
  supplier_id BIGINT,
  purchase_order_id BIGINT REFERENCES purchase_orders(id),
  product_id BIGINT,
  order_qty INTEGER,
  receive_order_qty INTEGER,
  stocked_order_qty INTEGER,
  price NUMERIC,
  total NUMERIC
);

CREATE TABLE receive_orders (
  id SERIAL PRIMARY KEY,
  supplier_id BIGINT,
  receive_date TIMESTAMP,
  status VARCHAR(50),
  tenant VARCHAR(100)
);

CREATE TABLE receive_order_items (
  id SERIAL PRIMARY KEY,
  receive_order_id BIGINT REFERENCES receive_orders(id),
  purchase_order_id BIGINT,
  product_id BIGINT,
  product_code VARCHAR(100),
  product_name VARCHAR(255),
  product_location VARCHAR(100),
  receive_order_qty INTEGER
);

CREATE TABLE sales_orders (
  id SERIAL PRIMARY KEY,
  customer_id BIGINT,
  order_date TIMESTAMP,
  amount NUMERIC,
  discount INTEGER,
  total NUMERIC,
  tenant VARCHAR(100)
);

CREATE TABLE sales_order_items (
  id SERIAL PRIMARY KEY,
  sales_order_id BIGINT REFERENCES sales_orders(id),
  product_id BIGINT,
  order_qty INTEGER,
  back_order_qty INTEGER,
  allocation_order_qty INTEGER,
  packing_order_qty INTEGER,
  invoice_order_qty INTEGER,
  price NUMERIC,
  total NUMERIC
);

CREATE TABLE suppliers (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(50),
  address TEXT,
  discount_percent BIGINT,
  term_of_payment VARCHAR(100),
  tenant VARCHAR(100)
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  password VARCHAR(255),
  role VARCHAR(50),
  tenant VARCHAR(100)
);
