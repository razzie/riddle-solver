package riddle

import (
	"encoding/json"
	"io/ioutil"
)

// Riddle contains the setup and rules of the riddle
type Riddle struct {
	Setup Setup
	Rules []Rule
}

// NewRiddle returns a new Riddle
func NewRiddle() *Riddle {
	return &Riddle{Setup: make(Setup)}
}

// LoadRiddle loads the riddle from a byte slice in JSON format
func LoadRiddle(data []byte) (*Riddle, error) {
	var r Riddle

	err := json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// LoadRiddleFromFile loads the riddle from a JSON file
func LoadRiddleFromFile(file string) (*Riddle, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return LoadRiddle(data)
}

// Save saves the riddle to a byte slice in JSON format
func (r *Riddle) Save() ([]byte, error) {
	return json.MarshalIndent(r, "", "\t")
}

// SaveToFile saves the riddle to a JSON file
func (r *Riddle) SaveToFile(file string) error {
	data, err := r.Save()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0777)
}
