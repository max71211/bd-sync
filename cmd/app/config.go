package main

type Config struct {
	AoptDB struct {
		Host     string `env:"AOPT_DB_HOST" envDefault:"localhost"`
		Port     string `env:"AOPT_DB_PORT" envDefault:"3306"`
		Database string `env:"AOPT_DB_DB" envDefault:"aopt"`
		User     string `env:"AOPT_DB_USER" envDefault:"root"`
		Password string `env:"AOPT_DB_PASSWORD" envDefault:""`
	}
	AutoDB struct {
		Host     string `env:"AUTO_DB_HOST" envDefault:"localhost"`
		Port     string `env:"AUTO_DB_PORT" envDefault:"3306"`
		Database string `env:"AUTO_DB_DB" envDefault:"auto"`
		User     string `env:"AUTO_DB_USER" envDefault:"root"`
		Password string `env:"AUTO_DB_PASSWORD" envDefault:""`
	}
}
