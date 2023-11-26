CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name text not null,
    phone_number text not null UNIQUE,
    otp text,
    otp_expiration_time timestamp
);