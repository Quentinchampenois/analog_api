package configs

type Config struct {
	Server   Server
	Database Database
}

func (c *Config) Load() {
	c.Server.Setup()
	c.Database.Setup()
}
