CREATE DATABASE tag_db;

CREATE TABLE customers (
	customer_id serial PRIMARY KEY,
	name VARCHAR(100) NOT NULL UNIQUE,
	customer_code VARCHAR(100)
);

CREATE TABLE warehouses (
	warehouse_id serial PRIMARY KEY,
	name VARCHAR(100) UNIQUE  NOT NULL
);

CREATE TABLE locations (
	location_id serial PRIMARY KEY,
	name VARCHAR(100)  NOT NULL,
	warehouse_id int REFERENCES warehouses(warehouse_id),
	CONSTRAINT locations_name_warehouse_id UNIQUE(name, warehouse_id)
);

CREATE TYPE material_type AS ENUM ('Carrier','Card','Envelope','Insert', 'Consumables');
CREATE TYPE owner AS ENUM('Tag', 'Customer');

CREATE TABLE materials (
	material_id serial,
	stock_id VARCHAR(100)  NOT NULL,
	location_id int REFERENCES locations(location_id),
	customer_id int REFERENCES customers(customer_id),
	material_type MATERIAL_TYPE  NOT NULL,
	description TEXT,
	notes TEXT,
	quantity int  NOT NULL,
	cost DECIMAL NOT NULL,
	min_required_quantity int,
	max_required_quantity int,
	updated_at TIMESTAMP,
	is_active BOOLEAN NOT NULL,
	owner OWNER NOT NULL,
	CONSTRAINT pk_location_stock_owner PRIMARY KEY (stock_id, location_id, owner)
);

CREATE TABLE transactions_log (
	transaction_id serial PRIMARY KEY,
	material_id int NOT NULL,
	stock_id VARCHAR(100) NOT NULL,
	quantity_change int NOT NULL,
	notes text,
	cost DECIMAL,
	job_ticket VARCHAR(100),
	updated_at timestamp,
	remaining_quantity int
);

CREATE TABLE incoming_materials (
	shipping_id SERIAL PRIMARY KEY,
	customer_id INT REFERENCES customers(customer_id),
	stock_id VARCHAR(100) NOT NULL,
	cost DECIMAL NOT NULL,
	quantity INT NOT NULL,
	min_required_quantity int,
	max_required_quantity int,
	notes VARCHAR(100),
	is_active BOOLEAN NOT NULL,
	type VARCHAR(100) NOT NULL,
	owner OWNER NOT NULL
);
