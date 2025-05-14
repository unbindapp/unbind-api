-- modify "templates" table
ALTER TABLE "templates" ADD COLUMN "description" character varying NOT NULL, ADD COLUMN "icon" character varying NOT NULL, ADD COLUMN "keywords" jsonb NULL;
