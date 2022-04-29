# hydra

[![Go CI](https://github.com/getpolygon/hydra/actions/workflows/go.yml/badge.svg)](https://github.com/getpolygon/hydra/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/getpolygon/hydra.svg)](https://pkg.go.dev/github.com/getpolygon/hydra)

Hydra is a hybrid configuration management library for Go, created for simplicity and testability. It supports using YAML and optionally paired with environment variables.

## Why?

In testing environments, you might not want to have separate configuration files for each one of your tests. Moreover, using more mature configuration libraries such as [Viper](https://github.com/spf13/viper), it is virtually impossible to write tests that use configuration files, since they have to be in a directory different than your project root, which honestly, makes stuff messy.

For this reason, we created Hydra, which will read, load, and fill in the blanks of the incomplete YAML configuration using environment variables. Think of it this way, if a YAML file does not exist, then Hydra will attempt to load the configuration using environment variables, optionally defined in your schema. However, if a configuration file was found, but has missing fields, Hydra will optionally fill in those fields with the values loaded from the environment.

## Installation

Get started by installing the latest version of hydra:

```bash
go get -u github.com/getpolygon/hydra
```

Define a configuration schema with validators and names for fallback environment variables:

```go
type Settings struct {
    Port        int16  `yaml:"port" env:"PORT"`
    Address     string `yaml:"address" env:"ADDRESS"`
    PostgreSQL  string `yaml:"postgres" validate:"uri" env:"DSN"`
}
```

create the configuration file:

`conf.yml`

```yml
port: 1234
address : "127.0.0.1"
postgres: "postgres://<postgres>:<password>@<host>:<port>/<database>?sslmode=disable"
```

and load the configuration by creating a `hydra.Hydra` instance:

```go
h := hydra.Hydra{
    Config: hydra.Config{
        Paths: []string{
            "~/somepath/config.yaml",
        },
    }
}

_, err := h.Load(new(Settings))
// ...
```

Initially, Hydra will attempt to find an existing configuration file provided from the paths in Hydra configuration, and then it will move on to filling in the missing values from environment variables.

## License

This software is licensed under the permissive [BSD-3-Clause license](./LICENSE).
