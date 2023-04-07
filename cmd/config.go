package cmd

type Config struct {
	Log struct {
		Level  string `json:"level"`
		Format string `json:"format"`
	} `json:"log"`

	Validator struct {
		Local     bool   `json:"local"`
		RulesPath string `json:"rulesPath"`

		DB struct {
			URI  string `json:"uri" yaml:"uri"`
			Name string `json:"name" yaml:"name"`
		} `json:"db" yaml:"db"`
	} `json:"validator"`
}
