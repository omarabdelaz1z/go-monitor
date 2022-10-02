package config

type Config struct {
	Log struct {
		Level string
		Path  string
	}

	Db struct {
		Driver       string
		Dsn          string
		MaxIdleConns int
		MaxOpenConns int
		MaxIdleTime  int
	}
}
