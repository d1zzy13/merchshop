package integration

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func SetupTestDB(t *testing.T) *sql.DB {

	testDB, err := sql.Open("postgres", "host=localhost user=dmitry password=mypass dbname=testdb sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	_, err = testDB.Exec(`
        CREATE TABLE users (
            id BIGSERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
            balance BIGINT NOT NULL CHECK (balance >= 0),
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE merchandise (
            name VARCHAR(50) PRIMARY KEY,
            price BIGINT NOT NULL CHECK (price > 0)
        );

        CREATE TABLE purchases (
            id BIGSERIAL PRIMARY KEY,
            user_id BIGINT NOT NULL REFERENCES users(id),
            merch_name VARCHAR(50) NOT NULL REFERENCES merchandise(name),
            quantity INT NOT NULL CHECK (quantity > 0),
            total_price BIGINT NOT NULL CHECK (total_price > 0),
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );
    `)

	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	_, err = testDB.Exec("TRUNCATE TABLE users, merchandise, purchases RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}

	_, err = testDB.Exec(`
        INSERT INTO users (username,password_hash, balance) VALUES ('alice','pass', 1000), ('bob','pass', 500);
        INSERT INTO merchandise (name, price) VALUES ('t-shirt', 80), ('cup', 20);
    `)

	if err != nil {
		t.Fatalf("failed to insert initial data: %v", err)
	}

	return testDB
}
