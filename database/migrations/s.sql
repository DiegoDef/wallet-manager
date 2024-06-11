CREATE TABLE cryptocurrency (
    cryptocurrency_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
	balance DECIMAL(30, 18),
	fiat_balance NUMERIC(14,2),
    created_date DATE NOT NULL DEFAULT CURRENT_DATE
);

CREATE TABLE crypto_transaction (
    transaction_id SERIAL PRIMARY KEY,
    cryptocurrency_id INT NOT NULL,
	cryptocurrency_amount NUMERIC(30, 18),
	fiat_amount NUMERIC(14,2), -- only dollar at the moment
	purchase_date DATE NOT NULL DEFAULT CURRENT_DATE,
	created_date DATE NOT NULL DEFAULT CURRENT_DATE,
    FOREIGN KEY (cryptocurrency_id) REFERENCES cryptocurrency (cryptocurrency_id) ON DELETE CASCADE
);

CREATE TABLE crypto_price (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price_usd NUMERIC(14,2) NOT NULL,
	updated_date DATE NOT NULL DEFAULT CURRENT_DATE
);

CREATE OR REPLACE FUNCTION get_percentage_profit(crypto_balance numeric, price_usd numeric, fiat_balance numeric)
returns NUMERIC
language plpgsql
as
$$
declare
begin
   RETURN (COALESCE(crypto_balance, 0) * COALESCE(PRICE_USD, 0)) / COALESCE(NULLIF(fiat_balance, 0), 1) * 100;
end;
$$;

CREATE OR REPLACE FUNCTION get_usd_profit(crypto_balance numeric, price_usd numeric, fiat_balance numeric)
returns numeric
language plpgsql
as
$$
declare
begin
   return (COALESCE(crypto_balance, 0) * COALESCE(PRICE_USD, 0)) - COALESCE(fiat_balance, 0);
end;
$$;