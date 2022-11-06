package main

type Config struct {
	Main struct {
		AppName             string `env:"APP_NAME" envDefault:"app"`
		UpdateBrands        bool   `env:"UPDATE_BRANDS" envDefault:"false"`
		UpdateVehicles      bool   `env:"UPDATE_VEHICLES" envDefault:"false"`
		UpdateModifications bool   `env:"UPDATE_MODIFICATIONS" envDefault:"false"`
	}
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
