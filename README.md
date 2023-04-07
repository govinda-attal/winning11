>Note: I may have gone above and beyond the expectations of exercise. But thought it is wise to make some assumptions and continue (as limited interaction in the begining and abstract nature of exercise). But helps team to judge better.

# Software Requirements

* Golang SDK 1.20.X installed
* Docker 23.0.X installed (only to run the application against MongoDB server)
* $GOPATH is set and $GOPATH/bin is added to $PATH

# Handy make commands

## Linting
```sh
make lint
```

## Tests
```sh
make test
```

## Build
```sh
make build
```

## Run
### Validation Rules from local file
```sh
make build
make run_local
```

### Validation Rules from external database
```sh
make build
make docker_up
make run_local
make docker_down
```

# Assumptions & Decisions

## Article and Feeds
Exercise description includes   
```
Given the below object:
{
“topic”: “A”,
“name”: “a”,
“description”: “something”
}

``` 
However it doesn't mention what this object represents. So made an assumption that this object is about an `Article` which includes following fields `topic, name & description`

In _*Go*_,  good package naming convention is driven by Domain Driven Design. Again made an assumption to name package as `feeds`. Thus feeds will provide functionality to work with an `Article` and its _*validations*_

## Validator
Exercise warrants to implement a *Validator* and its nature can be defined using interface
```go
type Validator interface {
	Validate(context.Context, Article) error
}
```
Validator validates an article and returns an error which lists *all errors* for an invalid one.

## Validation Rules Provider
Exercise mentioned:

```
Write a validator to check following rules:
1. If topic == A, then name will be “a”and description will be more than 10 and less than 100 chars
2. If topic == B, then name will be “x” and description will be less than 40 chars

Note: there will be more topic D, E, F …. In the future and combinations of validation rules on name and description are varying.
```

From this it clear that rules are different for each topic. However it is not exactly clear on number of topics and/or dynamic nature of validations.
- Small number of topics and static nature of validation rules
- Large (or unknown) number and/or dynamic nature of Topic and Validations

Therefore we have two options which can help to decide accordingly. Hence `RulesProvider` interface decouples core validation logic from rules fetcher/provider.

```go
type RulesProvider interface {
	Rules(ctx context.Context, topic string) (*ValidationRules, error)
	Shutdown(ctx context.Context) error
}
```

For now app supports:
- Local File Rules Provider
- External DB Rules Provider

Application configuration `app-cfg.yaml` enables to choose one of them at runtime.
[Sample Application Configuration](./configs/app-cfg.yaml)

```yaml
validator:
  local: true # if false provide valid db config
  rulesPath: ./configs/rules.yaml
  db:
    uri: mongodb://root:password@localhost:27017/?authSource=admin
    name: winning11
```

### Local File Rules Provider
Validation rules for these topics can be configured using local file. The rules configuration file can be in JSON, YAML, etc. For readability picked `YAML` format.
[Sample Validation Rules](./configs/rules.yaml)
```yaml
A:
  name:
    value: a
  desc:
    lenMoreThan: 10
    lenLessThan: 100

B:
  name: 
    value: b
  desc:
    lenLessThan: 40
```

### External Rules Provider
In such case validation rules for topics can be maintained in external (database, config-engine) enabling application to fetch this data from same. As Job Description mentioned NoSQL DB - MongoDB the app supports same. In theory this can be replaced by any database or other suitable providers.

* Migrations
[load sample article validations in articles collection](./migrations/1_article.go)


## Validation Logic
Validator logic is pretty simple:
```go
import validation "github.com/go-ozzo/ozzo-validation/v4"

func (v *validator) Validate(ctx context.Context, a Article) error {
	rules, err := v.rulesPrv.Rules(ctx, a.Topic)
	if err != nil {
		// unable to fetch rules for a given topic
		return err
	}

	return validation.ValidateStruct(&a,
		validation.Field(&a.Name,
			validation.By(
				valext.StrEquals(rules.Name.Value),
			),
		),
		validation.Field(&a.Description,
			validation.By(
				valext.StrLenBetween(rules.Desc.LenMoreThan, rules.Desc.LenLessThan),
			),
		),
	)
}
```

In my experience validation requirements can be very dynamic and complex in nature. It such cases third party libraries can be helpful. There are many options. 
For this example I have picked `github.com/go-ozzo/ozzo-validation`. It helps to implement dynamic custom validations. It provides a nice way to list *all errors* against one or more fields of a struct.

Validation rules of this exercise are slightly different than ones that are offerred out of the box from above library. Hence implemented custom rules in package [validations](./utils/validation/)
* StrLenBetween - checks if length of the string is between (X,Y) exclusive
* StrEquals - checks if given value equals a string

## Nature of Application
From exercise description it is not clear how to interact with the application. Is it a cli, webserivce, or library?

For simplicity have implemented it as a cli with multiple sub-commands
* validate: validate a given article (flag `--article=./sample-article.json`)
* migrate: helpful to load mongodb (migrate `up`) with validation rules

A sub-command `serve` can be included to run the application as a webservice providing an endpoint to validate an article.

### Local Execution against validation rules in a file
```sh
make build
# view/change dist/configs/rules.yaml and/or dist/sample-article.json
make run_local
```