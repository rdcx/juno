package mysql

import "database/sql"

var migrations = map[string]string{
	"create_queries_table": `
		CREATE TABLE IF NOT EXISTS queries (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			status VARCHAR(16) NOT NULL DEFAULT 'pending',
			query_type VARCHAR(16) NOT NULL,
			basic_query_version VARCHAR(16) NOT NULL,
			basic_query MEDIUMTEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);`,
}

func ExecuteMigrations(db *sql.DB) error {

	// create migrations table if it doesn't exist
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
			name VARCHAR(255) NOT NULL PRIMARY KEY
		);`); err != nil {
		return err
	}

	for name, migration := range migrations {

		// check if migration has already been executed
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", name).Scan(&count); err != nil {
			return err
		}

		if count > 0 {
			continue
		}

		if _, err := db.Exec(migration); err != nil {
			return err
		}

		db.Exec("INSERT INTO migrations (name) VALUES (?)", name)
	}

	return nil
}
