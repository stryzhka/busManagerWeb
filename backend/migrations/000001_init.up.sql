CREATE TABLE "drivers" (
                           "id"	TEXT UNIQUE,
                           "name"	TEXT NOT NULL,
                           "surname"	TEXT NOT NULL,
                           "patronymic"	TEXT NOT NULL,
                           "birth_date"	TIMESTAMP NOT NULL,
                           "passport_series"	TEXT NOT NULL UNIQUE,
                           "snils"	TEXT NOT NULL UNIQUE,
                           "license_series"	TEXT NOT NULL UNIQUE,
                           PRIMARY KEY("id")
);

CREATE TABLE "buses" (
                         "id"	TEXT UNIQUE,
                         "brand"	TEXT NOT NULL,
                         "bus_model"	TEXT NOT NULL,
                         "register_number"	TEXT NOT NULL UNIQUE,
                         "assembly_date"	TIMESTAMP NOT NULL,
                         "last_repair_date"	TIMESTAMP NOT NULL,
                         PRIMARY KEY("id")
);

CREATE TABLE "routes" (
                          "id"	TEXT UNIQUE,
                          "number"	TEXT UNIQUE,
                          PRIMARY KEY("id")
);

CREATE TABLE "bus_stops" (
                             "id"	TEXT NOT NULL UNIQUE,
                             "lat"	REAL NOT NULL UNIQUE,
                             "long"	REAL NOT NULL UNIQUE,
                             "name"	TEXT NOT NULL UNIQUE,
                             PRIMARY KEY("id")
);

CREATE TABLE "routes_drivers" (
                                  "route_id"	TEXT,
                                  "driver_id"	TEXT
);

CREATE TABLE "routes_buses" (
                                "route_id"	TEXT,
                                "bus_id"	TEXT
);

CREATE TABLE "routes_bus_stops" (
                                    "route_id"	TEXT,
                                    "bus_stop_id"	TEXT
);