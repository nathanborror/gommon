package settings

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Database represents database settings
type Database struct {
	Kind string
	Name string
	User string
}

// Settings represents common settings
type Settings struct {
	Database Database
}

// NewSettings returns a new settings struct
func NewSettings(path string) Settings {
	file, _ := os.Open(path)
	decoder := json.NewDecoder(file)
	settings := Settings{}

	err := decoder.Decode(&settings)
	if err != nil {
		log.Println(err)
	}

	return settings
}

// DataSource returns dataSource string for use with sql.Connect
func (s Settings) DataSource() string {
	if s.Database.Kind == "sqlite3" {
		return s.Database.Name
	}
	return fmt.Sprintf("user=%s dbname=%s sslmode=disable", s.Database.User, s.Database.Name)
}

// DriverName returns the database driver name for use with sql.Connect
func (s Settings) DriverName() string {
	return s.Database.Kind
}
