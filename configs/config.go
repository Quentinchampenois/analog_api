package configs

type Config struct {
	Server      Server
	Database    Database
	Application Application
}

func (c *Config) Load() {
	c.Server.Setup()
	c.Database.Setup()
	c.Application.Setup()
}
