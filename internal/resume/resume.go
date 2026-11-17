package resume

import (
	"fmt"
	"os"
	"time"

	yaml "github.com/goccy/go-yaml"
)

type Resume struct {
	Name       Name         `yaml:"name"`
	Contact    Contact      `yaml:"contact"`
	Address    Address      `yaml:"address"`
	Experience []Experience `yaml:"experience"`
	Projects   []Projects   `yaml:"projects"`
}

type Name struct {
	FirstName  string `yaml:"first_name"`
	MiddleName string `yaml:"middle_name,omitempty"`
	LastName   string `yaml:"last_name"`
}

type Contact struct {
	EmailAddress string `yaml:"email_address"`
	PhoneNumber  string `yaml:"phone_number"`
	FaxNumber    string `yaml:"fax_number"`
}

type Address struct {
	HomeAddress       string   `yaml:"home_address"`
	WorkableRemote    bool     `yaml:"workable_remote"`
	WorkableInOffice  bool     `yaml:"workable_in_office"`
	WorkableStates    []string `yaml:"workable_states"`
	WorkableAllCities bool     `yaml:"workable_all_cities"`
	WorkableCities    []string `yaml:"workable_cities"`
	WorkableCommute   string   `yaml:"workable_commute"`
}

type Experience struct {
	CompanyName string    `yaml:"company_name"`
	JobTitle    string    `yaml:"job_title"`
	Location    Location  `yaml:"location"`
	StartDate   time.Time `yaml:"start_date,omitempty"`
	EndDate     time.Time `yaml:"end_date,omitempty"`
	Bullets     []Bullet  `yaml:"bullets"`
}

type Bullet struct {
	Text     string   `yaml:"text"`
	Keywords []string `yaml:"keywords"`
}

type Location struct {
	State string `yaml:"state"`
	City  string `yaml:"city"`
}

type Links struct {
	LiveUrl   string `yaml:"live_url"`
	GitHubUrl string `yaml:"github_url"`
}

type Projects struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Links       Links     `yaml:"links"`
	StartDate   time.Time `yaml:"start_date,omitempty"`
	EndDate     time.Time `yaml:"end_date,omitempty"`
	Bullets     []Bullet  `yaml:"bullets"`
}

// Decode .yaml path file into resume struct
func (r *Resume) Load(p string) error {

	// TODO: Move file path handling to dedicated package
	info, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("getting file %s: %w", p, err)
	}
	if info.IsDir() {
		return fmt.Errorf("is directory %s: %w", p, err)
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, r); err != nil {
		return err
	}

	return nil
}
