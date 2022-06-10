CREATE TABLE restaurants (id SERIAL PRIMARY KEY,name char(50) NOT NULL,avg_time interval NOT NULL,avg_price integer NOT NULL);
CREATE TABLE tables (id SERIAL PRIMARY KEY,restaurant_id integer REFERENCES restaurants (id),capacity integer NOT NULL);
CREATE TABLE clients (id SERIAL PRIMARY KEY,name char(50) NOT NULL,phone_number char(50) NOT NULL, UNIQUE (name, phone_number));
CREATE TABLE reservations (id SERIAL PRIMARY KEY,table_id integer REFERENCES tables (id),start_time timestamp NOT NULL,end_time timestamp NOT NULL,reserved_by int REFERENCES clients (id));