CREATE TABLE IF NOT EXISTS auth$users
		(
			name VARCHAR(255) NOT NULL UNIQUE,
		    password VARCHAR(255) NOT NULL
		);

CREATE TABLE IF NOT EXISTS products
		(
		    id SERIAL PRIMARY KEY,
		    name VARCHAR(255) NOT NULL,
		    price INT NOT NULL,
		    quantity INT NOT NULL,
		    owner_name VARCHAR(255) NOT NULL,
		    created_at TIMESTAMP NOT NULL,
		    FOREIGN KEY (owner_name) REFERENCES auth$users(name)
		);
