CREATE TABLE IF NOT EXISTS "fields" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "name" TEXT(100) NOT NULL,
  "lat" REAL(10,8),
  "lng" REAL(10,8),
  "description" TEXT,
  "image_name" TEXT(255),
  "image_original_name" TEXT(255),
  "image_mime_type" TEXT(255),
  "image_size" INTEGER(11),
  "created_at" TEXT(20) NOT NULL,
  "updated_at" TEXT(20)
);

CREATE TABLE IF NOT EXISTS "reservoirs" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "field_id" INTEGER(11),
  "name" TEXT(100) NOT NULL,
  "capacity" REAL(10, 2) NOT NULL,
  "measurement_unit" INTEGER(11),
  "created_at" TEXT(20) NOT NULL,
  "updated_at" TEXT(20),
  CONSTRAINT "reservoirs_field_id_foreign" FOREIGN KEY ("field_id") REFERENCES "fields" ("id")
);
CREATE INDEX IF NOT EXISTS reservoirs_field_id ON reservoirs (field_id);

CREATE TABLE IF NOT EXISTS "areas" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "reservoir_id" INTEGER(11),
  "field_id" INTEGER(11),
  "name" TEXT(100) NOT NULL,
  "growing_method" INTEGER(11) NOT NULL,
  "capacity" INTEGER(11) NOT NULL,
  "measurement_unit" INTEGER(11) NOT NULL,
  "image_name" TEXT(255),
  "image_original_name" TEXT(255),
  "image_mime_type" TEXT(255),
  "image_size" INTEGER(11),
  "created_at" TEXT(20) NOT NULL,
  "updated_at" TEXT(20),
  CONSTRAINT "areas_field_id_foreign" FOREIGN KEY ("field_id") REFERENCES "fields" ("id"),
  CONSTRAINT "areas_reservoir_id_foreign" FOREIGN KEY ("reservoir_id") REFERENCES "reservoirs" ("id")
);
CREATE INDEX IF NOT EXISTS areas_field_id ON areas (field_id);
CREATE INDEX IF NOT EXISTS areas_reservoir_id ON areas (reservoir_id);

CREATE TABLE IF NOT EXISTS "devices" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "field_id" INTEGER(11),
  "name" TEXT(100),
  "description" TEXT,
  "device_type" INTEGER(11),
  "created_at" TEXT(20) NOT NULL,
  "updated_at" TEXT(20),
  CONSTRAINT "devices_field_id_foreign" FOREIGN KEY ("field_id") REFERENCES "fields" ("id")
);
CREATE INDEX IF NOT EXISTS devices_field_id ON devices (field_id);

CREATE TABLE IF NOT EXISTS "areas_devices" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "area_id" INTEGER(11),
  "device_id" INTEGER(11),
  "created_at" TEXT(20) NOT NULL,
  "updated_at" TEXT(20),
  CONSTRAINT "areas_devices_areaa_id_foreign" FOREIGN KEY ("area_id") REFERENCES "areas" ("id"),
  CONSTRAINT "areas_devices_device_id_foreign" FOREIGN KEY ("device_id") REFERENCES "devices" ("id")
);
CREATE INDEX IF NOT EXISTS areas_devices_area_id ON areas_devices (area_id);
CREATE INDEX IF NOT EXISTS areas_devices_device_id ON areas_devices (device_id);

CREATE TABLE IF NOT EXISTS "seed_categories" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "name" TEXT(100) NOT NULL,
  "slug" TEXT(100) NOT NULL,
  "description" TEXT,
  "created_at" TEXT(20) NOT NULL,
  "updated_at" TEXT(20)
);

CREATE TABLE IF NOT EXISTS "seeds" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "seedcategory_id" INTEGER(11),
  "name" TEXT(100) NOT NULL,
  "quantity" INTEGER(11) NOT NULL,
  "measurement_unit" INTEGER(11) NOT NULL,
  "producer_name" TEXT(150) NOT NULL,
  "origin_country" TEXT(100) NOT NULL,
  "note" longtext,
  "expiration_month" TEXT(20) NOT NULL,
  "expiration_year" INTEGER(11) NOT NULL,
  "germination_rate" decimal(5,2),
  "image_name" TEXT(255),
  "image_original_name" TEXT(255),
  "image_mime_type" TEXT(255),
  "image_size" INTEGER(11),
  "updated_at" TEXT(20),
  "created_at" TEXT(20) NOT NULL,
  CONSTRAINT "seeds_seedcategory_id_foreign" FOREIGN KEY ("seedcategory_id") REFERENCES "seed_categories" ("id")
);
CREATE INDEX IF NOT EXISTS seeds_seedcategory_id ON seeds (seedcategory_id);

CREATE TABLE IF NOT EXISTS "plants" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "area_id" INTEGER(11) NOT NULL,
  "seed_id" INTEGER(11) NOT NULL,
  "seedling_date" TEXT(10),
  "area_capacity" INTEGER(11) NOT NULL,
  "harvesting_date" TEXT(10),
  "disposing_date" TEXT(10),
  "note" TEXT,
  "action" TEXT(10),
  "created_at" TEXT(20) NOT NULL,
  "updated_at" TEXT(20),
  CONSTRAINT "plants_area_id_foreign" FOREIGN KEY ("area_id") REFERENCES "areas" ("id"),
  CONSTRAINT "plants_seed_id_foreign" FOREIGN KEY ("seed_id") REFERENCES "seeds" ("id")
);
CREATE INDEX IF NOT EXISTS plants_area_id ON plants (area_id);
CREATE INDEX IF NOT EXISTS plants_seed_id ON plants (seed_id);

CREATE TABLE IF NOT EXISTS "resources" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "type" TEXT(100) NOT NULL,
  "created_at" TEXT(20) NOT NULL,
  "updated_at" TEXT(20)
);

CREATE TABLE IF NOT EXISTS "resources_devices" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "device_id" INTEGER(11),
  "resource_id" INTEGER(11),
  "name" TEXT(100) NOT NULL,
  "description" TEXT,
  "rid" TEXT(100) NOT NULL,
  "data_type" TEXT(20) NOT NULL,
  "unit" TEXT(20) NOT NULL,
  "created_at" TEXT(20) NOT NULL,
  "updated_at" TEXT(20),
  CONSTRAINT "resources_devices_device_id_foreign" FOREIGN KEY ("device_id") REFERENCES "devices" ("id"),
  CONSTRAINT "resources_devices_resource_id_foreign" FOREIGN KEY ("resource_id") REFERENCES "resources" ("id")
);
CREATE INDEX IF NOT EXISTS resources_devices_area_id ON resources_devices (device_id);
CREATE INDEX IF NOT EXISTS resources_devices_seed_id ON resources_devices (resource_id);

CREATE TABLE IF NOT EXISTS "settings" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "key" TEXT(200) NOT NULL,
  "value" TEXT(200) NOT NULL,
  "updated_at" TEXT(20)
);

CREATE TABLE IF NOT EXISTS "tasks" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "name" TEXT(100) NOT NULL,
  "notes" TEXT,
  "category" TEXT(50) NOT NULL,
  "due_date" TEXT(20) NOT NULL,
  "urgency_level" TEXT(15) NOT NULL,
  "is_done" INTEGER(11),
  "field_id" INTEGER(11),
  "created_at" TEXT(20) NOT NULL,
  "updated_at" TEXT(20)
);

CREATE TABLE IF NOT EXISTS "users" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "username" TEXT(180) NOT NULL,
  "username_canonical" TEXT(180) NOT NULL,
  "email" TEXT(180) NOT NULL,
  "email_canonical" TEXT(180) NOT NULL,
  "enabled" INTEGER(1) NOT NULL,
  "salt" TEXT(255),
  "password" TEXT(255) NOT NULL,
  "last_login" TEXT(20),
  "confirmation_token" TEXT(180),
  "password_requested_at" TEXT(20),
  "roles" TEXT NOT NULL,
  "created_at" TEXT(20) NOT NULL,
  "updated_at" TEXT(20)
);
