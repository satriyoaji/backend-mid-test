
CREATE TABLE "employees" (
     "id" serial not null,
     "first_name" varchar not null,
     "last_name" varchar not null,
     "email" varchar unique not null,
     "hire_date" timestamptz not null,
     "created_at" timestamptz not null default current_timestamp,
     "updated_at" timestamptz not null default current_timestamp
);
