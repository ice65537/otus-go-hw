DROP SCHEMA IF EXISTS clndr CASCADE;
DROP ROLE IF EXISTS clndr;
CREATE ROLE clndr WITH
  LOGIN
  NOSUPERUSER
  INHERIT
  NOCREATEDB
  NOCREATEROLE
  NOREPLICATION
  PASSWORD 'clndr';

CREATE SCHEMA IF NOT EXISTS clndr AUTHORIZATION clndr;
GRANT ALL ON SCHEMA clndr TO clndr;

DROP TABLE IF EXISTS clndr.t_event CASCADE;
CREATE TABLE IF NOT EXISTS clndr.t_event
(
    eid character varying(100) COLLATE pg_catalog."default" NOT NULL,
    etitle character varying(300) COLLATE pg_catalog."default" NOT NULL,
    estartdt timestamp with time zone NOT NULL,
    estopdt timestamp with time zone NOT NULL,
    edesc text COLLATE pg_catalog."default",
    eowner character varying(100) COLLATE pg_catalog."default" NOT NULL,
    enotifybefore integer
)
TABLESPACE pg_default;
ALTER TABLE IF EXISTS clndr.t_event OWNER to clndr;

DROP INDEX IF EXISTS clndr.idx_event_estartdt;
CREATE INDEX IF NOT EXISTS idx_event_estartdt
    ON clndr.t_event USING btree
    (estartdt ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS clndr.uq_event_eid;
CREATE UNIQUE INDEX IF NOT EXISTS uq_event_eid
    ON clndr.t_event USING btree
    (eid COLLATE pg_catalog."default" varchar_ops ASC NULLS LAST)
    TABLESPACE pg_default;