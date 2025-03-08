package configs

import (
	"github.com/pkg/errors"
)

const (
	test        EnvironmentName = "test"
	development EnvironmentName = "development"
	staging     EnvironmentName = "staging"
	production  EnvironmentName = "production"
)

type EnvironmentName string

type Environment struct {
	name EnvironmentName
}

func NewEnvironmentFromRawEnvVar(name string) Environment {
	guardEnvironmentName(name)
	return Environment{name: EnvironmentName(name)}
}

func environments() map[EnvironmentName]struct{} {
	return map[EnvironmentName]struct{}{
		test:        {},
		development: {},
		staging:     {},
		production:  {},
	}
}

func guardEnvironmentName(name string) {
	env := EnvironmentName(name)

	if _, environmentExists := environments()[env]; !environmentExists {
		panic(errors.Errorf("environment <%s> doesnt exist", name))
	}
}

func (env Environment) IsTest() bool {
	return env.name == test
}

func (env Environment) IsDevelopment() bool {
	return env.name == development
}

func (env Environment) IsStaging() bool {
	return env.name == staging
}

func (env Environment) IsProduction() bool {
	return env.name == production
}
