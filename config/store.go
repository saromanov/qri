package config

import "github.com/qri-io/jsonschema"

// Store configures a qri content addessed file store (cafs)
type Store struct {
	Type string `json:"type"`
}

// DefaultStore returns a new default Store configuration
func DefaultStore() *Store {
	return &Store{
		Type: "ipfs",
	}
}

// Validate validates all fields of store returning all errors found.
func (cfg Store) Validate() error {
	schema := jsonschema.Must(`{
    "$schema": "http://json-schema.org/draft-06/schema#",
    "title": "Store",
    "description": "Config for the qri content addressed file store",
    "type": "object",
    "required": ["type"],
    "properties": {
      "type": {
        "description": "Type of store",
        "type": "string",
        "enum": [
          "ipfs"
        ]
      }
    }
  }`)
	return validate(schema, &cfg)
}

// Copy returns a deep copy of the Store struct
func (cfg *Store) Copy() *Store {
	res := &Store{
		Type: cfg.Type,
	}

	return res
}