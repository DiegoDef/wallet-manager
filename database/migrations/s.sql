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