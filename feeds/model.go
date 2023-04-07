package feeds

type Article struct {
	Topic       string `json:"topic"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ArticleValidationRules map[string]ValidationRules

type ValidationRules struct {
	Name ValueRule `json:"name" yaml:"name"`
	Desc LenRule   `json:"desc" yaml:"desc"`
}

type ValueRule struct {
	Value string `json:"value" yaml:"value"`
}

type LenRule struct {
	LenMoreThan int `json:"lenMoreThan" yaml:"lenMoreThan"`
	LenLessThan int `json:"lenLessThan" yaml:"lenLessThan"`
}

type ArticleValidation struct {
	Topic           string          `bson:"topic" json:"topic"`
	ValidationRules ValidationRules `bson:"validationRules" json:"validationRules"`
}
