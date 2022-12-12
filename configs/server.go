package configs

import (
	"fmt"
	"log"
	"os"
)

type Server struct {
	Port      string
	Host      string
	Ssl       string
	JWTSecret []byte
}

func (s *Server) Setup() {
	s.Port = s.getPort()
	s.Host = s.getHost()
	s.Ssl = s.getHttp()
	s.JWTSecret = s.getJWTSecret()
}

func (s *Server) GetFullPath() string {
	return fmt.Sprintf("%v://%v:%v/", s.Ssl, s.Host, s.Port)
}

func (s *Server) getPort() string {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		return port
	}

	return "8080"
}

func (s *Server) getHost() string {
	if host := os.Getenv("SERVER_HOST"); host != "" {
		return host
	}

	return "localhost"
}

func (s *Server) getHttp() string {
	if http := os.Getenv("SERVER_SSL"); http != "" {
		return "https"
	}

	return "http"
}

func (s *Server) getJWTSecret() []byte {
	if jwtSecret := os.Getenv("JWT_SECRET_PASSWORD"); jwtSecret != "" {
		return []byte(jwtSecret)
	}

	log.Fatalln("You must define 'JWT_SECRET_PASSWORD' environment variable for JWT authentication system")
	return nil
}
