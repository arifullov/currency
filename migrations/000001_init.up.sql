CREATE EXTENSION citext;

CREATE TABLE IF NOT EXISTS currencies (
  currency_from citext NOT NULL,
  currency_to citext NOT NULL,
  well numeric NOT NULL,
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  PRIMARY KEY (currency_from, currency_to)
);
