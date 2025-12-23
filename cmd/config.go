package cmd

import "os"

type Config struct {
	ServerAddress string
	ListenPort    string
	EncryptionKey []byte
	MTU           int
}

func LoadConfig() Config {
	return Config{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		EncryptionKey: []byte(os.Getenv("ENCRYPTION_KEY")),
		ListenPort:    os.Getenv("LISTEN_PORT"),
		MTU:           1500,
	}
}
