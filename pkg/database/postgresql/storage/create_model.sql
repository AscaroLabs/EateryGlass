CREATE TABLE restaurants (id SERIAL PRIMARY KEY,name TEXT NOT NULL,avg_time interval NOT NULL,avg_price integer NOT NULL, UNIQUE(name));
CREATE TABLE tables (id SERIAL PRIMARY KEY,restaurant_id integer REFERENCES restaurants (id),capacity integer NOT NULL);
CREATE TABLE clients (id SERIAL PRIMARY KEY,name TEXT NOT NULL,phone_number TEXT NOT NULL, UNIQUE (name, phone_number));
CREATE TABLE reservations (id SERIAL PRIMARY KEY,table_id integer REFERENCES tables (id),start_time timestamp NOT NULL,end_time timestamp NOT NULL,reserved_by int REFERENCES clients (id), UNIQUE (table_id, start_time));