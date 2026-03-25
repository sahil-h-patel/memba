package auth

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func GetToken() ([]byte, error) {
	// Get home dir
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	// get file path
	file := filepath.Join(home, ".memba", "token")

	// read file back
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func SaveToken(c *gin.Context) error {
	// Get token
	token, err := c.Cookie("Auth")
	if err != nil {
		c.String(500, "Failed to retrieve token from session")
		return err
	}

	// get home dir
	home, error := os.UserHomeDir()
	if error != nil {
		return error
	}

	dir := filepath.Join(home, ".memba")
	path := filepath.Join(dir, "token")

	// create dir if doesnt exist
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}

	// write token
	return os.WriteFile(path, []byte(token), 0600)
}
