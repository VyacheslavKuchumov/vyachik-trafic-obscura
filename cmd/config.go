package cmd

import "os"

type Config struct {
	ServerAddress string
	ListenPort    int
	EncryptionKey []byte
	MTU           int
}

func LoadConfig() Config {
	return Config{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		EncryptionKey: []byte(os.Getenv("ENCRYPTION_KEY")),
		MTU:           1500,
	}
}
