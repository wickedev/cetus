package main

type Image struct {
	Dockerfile string
}

type Label = map[string]string
type Rule = map[string]string
type Environment = map[string]string
type Port struct {
	Port string
}

//noinspection GoSnakeCaseUsage
type Service struct {
	Chart      Chart    `yaml:",omitempty"`
	Depends_on []string `yaml:",omitempty"`

	Image       string      `yaml:",omitempty"`
	Labels      []Label     `yaml:",omitempty"`
	Ingress     Ingress     `yaml:",omitempty"`
	Ports       []Port      `yaml:",omitempty"`
	Environment Environment `yaml:",omitempty"`
	Deploy      Deploy      `yaml:",omitempty"`
}

type Ingress struct {
	TLS   bool   `yaml:",omitempty"`
	Rules []Rule `yaml:",omitempty"`
}

type Deploy struct {
	Strategy   string      `yaml:",omitempty"` // ab, red-green, canary
	Variations []Variation `yaml:",omitempty"`
}

type Variation struct {
	Labels []Label `yaml:",omitempty"`
	Weight int     `yaml:",omitempty"`
}

type Chart struct {
	Name    string            `yaml:",omitempty"`
	Version string            `yaml:",omitempty"`
	Values  map[string]string `yaml:",omitempty"`
}

//noinspection GoSnakeCaseUsage
type Config struct {
	Version     string             `yaml:",omitempty"`
	Name        string             `yaml:",omitempty"`
	App_version string             `yaml:",omitempty"`
	Images      map[string]Image   `yaml:",omitempty"`
	Services    map[string]Service `yaml:",omitempty"`
}
