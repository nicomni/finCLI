package cli

import (
	"fincli/internal/bankcsv/autofill"
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

type Config interface {
	GetAutofillRules() ([]autofill.Rule, error) // Method to get autofill rules
	SetConfigFile(cfgFile string)               // Method to set the configuration file
}

type defaultConfig struct {
	v                 *viper.Viper
	hasReadConfigFile bool // Flag to indicate if the config file has been read
}

func New() Config {
	v := viper.New()
	return &defaultConfig{v: v}
}

func (d *defaultConfig) GetAutofillRules() ([]autofill.Rule, error) {
	if !d.hasReadConfigFile {
		if err := d.v.ReadInConfig(); err != nil {
			return nil, err // Return error if config file reading fails
		}
		d.hasReadConfigFile = true // Set the flag to true after reading the config file
	}
	var rules []autofill.Rule
	opt := viper.DecodeHook(toConditionTypeAndValueHookFunc)
	err := d.v.UnmarshalKey("autofill.rules", &rules, opt)
	return rules, err
}

func (d *defaultConfig) SetConfigFile(cfgFile string) {
	d.v.SetConfigFile(cfgFile)
}

func (d *defaultConfig) ReadInConfig() error {
	return d.v.ReadInConfig() // Reads the configuration file
}

var _ Config = &defaultConfig{}

func toConditionTypeAndValueHookFunc(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeOf(autofill.Condition{}) {
		return data, nil
	}
	dataVal, ok := data.(map[string]any)
	if !ok {
		return data, fmt.Errorf("expected map[string]any, got %T", data)
	}
	mapVal, ok := dataVal["equals"]
	if !ok {
		return data, fmt.Errorf("expected map to contain 'equals' key, got %v", dataVal)
	}
	strVal, ok := mapVal.(string)
	if !ok {
		return data, fmt.Errorf("expected value for 'equals' to be a string, got %T", mapVal)
	}

	delete(dataVal, "equals")
	dataVal["type"] = "equals"
	dataVal["value"] = strVal

	return dataVal, nil // Return the modified data
}
