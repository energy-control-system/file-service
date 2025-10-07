package config

type Settings struct {
	Port      int       `json:"port"`
	Databases Databases `json:"databases"`
}

type Databases struct {
	Minio    Minio  `json:"minio"`
	Postgres string `json:"postgres"`
}

type Minio struct {
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Password string `json:"password"`
	UseSSL   bool   `json:"use_ssl"`
	Host     string `json:"host"`
}
