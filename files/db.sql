
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS "customers";
CREATE TABLE "customers"(
  customer_id varchar(64) DEFAULT uuid_generate_v4() not null,
  customer_name varchar(80) not null,
  email varchar(50) not null,
  phone_number varchar(20) not null,
  dob timestamp(6) not null,
  sex boolean not null,
  password varchar(400) not null,
  "created_at" timestamp(6) DEFAULT now(),
  "updated_at" timestamp(6) DEFAULT now(),
  "deleted_at" timestamp(6),
  PRIMARY KEY(customer_id)
);

DROP TABLE IF EXISTS "products";
CREATE TABLE "products"(
  product_id varchar(64) DEFAULT uuid_generate_v4() not null,
  product_name varchar(80),
  basic_price numeric,
  "created_at" timestamp(6) DEFAULT now(),
  "updated_at" timestamp(6) DEFAULT now(),
  "deleted_at" timestamp(6),
  PRIMARY KEY(product_id)
);

DROP TABLE IF EXISTS "payment_methods";
CREATE TABLE "payment_methods"(
  payment_method_id varchar(64) DEFAULT uuid_generate_v4() not null,
  method_name varchar(70),
  code varchar(10) not null,
  "created_at" timestamp(6) DEFAULT now(),
  "updated_at" timestamp(6) DEFAULT now(),
  "deleted_at" timestamp(6),
  PRIMARY KEY(payment_method_id)
);

DROP TABLE IF EXISTS "orders";
CREATE TABLE "orders"(
  order_id varchar(64) DEFAULT uuid_generate_v4() not null,
  customer_id varchar(64) not null,
  order_number varchar(40) not null,
  order_date timestamp(6) not null,
  payment_method_id varchar(64) not null,
  "created_at" timestamp(6) DEFAULT now(),
  "updated_at" timestamp(6) DEFAULT now(),
  "deleted_at" timestamp(6),
  PRIMARY KEY(order_id),
  CONSTRAINT fk_customer FOREIGN KEY(customer_id) REFERENCES customers(customer_id),
  CONSTRAINT fk_payment_methods FOREIGN KEY(payment_method_id) REFERENCES payment_methods(payment_method_id)
);


DROP TABLE IF EXISTS "order_details";
CREATE TABLE "order_details"(
  order_detail_id varchar(64) DEFAULT uuid_generate_v4() not null,
  order_id varchar(64) not null,
  product_id varchar(64),
  qty int,
  "created_at" timestamp(6) DEFAULT now(),
  "updated_at" timestamp(6) DEFAULT now(),
  "deleted_at" timestamp(6),
  PRIMARY KEY(order_detail_id),
  CONSTRAINT fk_order FOREIGN KEY(order_id) REFERENCES orders(order_id),
  CONSTRAINT fk_product FOREIGN KEY(product_id) REFERENCES products(product_id)
);
