package mysql

import "database/sql"

var migrations = map[string]string{
	"create_strategies_table": `
		CREATE TABLE IF NOT EXISTS strategies (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);`,
	"create_strategy_fields_table": `
		CREATE TABLE IF NOT EXISTS strategy_fields (
			strategy_id VARCHAR(36) NOT NULL,
			field_id VARCHAR(36) NOT NULL,
			PRIMARY KEY (strategy_id, field_id)
		);`,
	"create_strategy_selectors_table": `
		CREATE TABLE IF NOT EXISTS strategy_selectors (
			strategy_id VARCHAR(36) NOT NULL,
			selector_id VARCHAR(36) NOT NULL,
			PRIMARY KEY (strategy_id, selector_id)
		);`,
	"create_strategy_filters_table": `
		CREATE TABLE IF NOT EXISTS strategy_filters (
			strategy_id VARCHAR(36) NOT NULL,
			filter_id VARCHAR(36) NOT NULL,
			PRIMARY KEY (strategy_id, filter_id)
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
