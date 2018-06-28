package particlemsg

import "os"

// GetServerConfig - returns these environment variables: PMSG_HOST, PMSG_PORT, PMSG_SSL_CERT, PMSG_SSL_KEY
func GetServerConfig() (string, string, string, string) {
	return os.Getenv("PMSG_HOST"), os.Getenv("PMSG_PORT"), os.Getenv("PMSG_SSL_CERT"), os.Getenv("PMSG_SSL_KEY")
}

// GetClientConfig - returns these environment variables: PMSG_NAME, PMSG_SSL_CERT, PMSG_SSL_KEY
func GetClientConfig() (string, string, string) {
	return os.Getenv("PMSG_NAME"), os.Getenv("PMSG_SSL_CERT"), os.Getenv("PMSG_SSL_KEY")
}
