package config

import "os"

// Config holds database connection strings and other configurations.
type Config struct {
	NodeDB          string
	UserDB          string
	BalancerDB      string
	TranDB          string
	ExtractionJobDB string
	SelectorDB      string
	FilterDB        string
	FieldDB         string
	StrategyDB      string
	RanagDB         string
}

// LoadConfig reads environment variables and returns a Config struct.
func LoadConfig() *Config {
	return &Config{
		NodeDB:          getEnv("NODE_DB", "root:juno@tcp(localhost:3306)/node?parseTime=true"),
		UserDB:          getEnv("USER_DB", "root:juno@tcp(localhost:3306)/user?parseTime=true"),
		BalancerDB:      getEnv("BALANCER_DB", "root:juno@tcp(localhost:3306)/balancer?parseTime=true"),
		TranDB:          getEnv("TRAN_DB", "root:juno@tcp(localhost:3306)/transaction?parseTime=true"),
		ExtractionJobDB: getEnv("EXTRACTION_JOB_DB", "root:juno@tcp(localhost:3306)/job?parseTime=true"),
		SelectorDB:      getEnv("SELECTOR_DB", "root:juno@tcp(localhost:3306)/selector?parseTime=true"),
		FilterDB:        getEnv("FILTER_DB", "root:juno@tcp(localhost:3306)/filter?parseTime=true"),
		FieldDB:         getEnv("FIELD_DB", "root:juno@tcp(localhost:3306)/field?parseTime=true"),
		StrategyDB:      getEnv("STRATEGY_DB", "root:juno@tcp(localhost:3306)/strategy?parseTime=true"),
		RanagDB:         getEnv("RANAG_DB", "root:juno@tcp(localhost:3306)/ranag?parseTime=true"),
	}
}

// getEnv reads an environment variable or returns a default value if not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
