CREATE TABLE IF NOT EXISTS "users" (
    "id" serial PRIMARY KEY NOT NULL,
    "first_name" varchar(50) NOT NULL,
    "email" varchar(100) not null,
    "created_at" timestamp DEFAULT current_timestamp
);