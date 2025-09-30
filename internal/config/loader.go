package config

import (
	"errors"
	"github.com/goloop/env"
)

func (c *Config) LoadYaml() error {
	panic("метод загрузки Yaml конфига пока не реализован")
}

func (c *Config) LoadEnv() error {
	if env.Load(".env") != nil {
		return errors.New("не удалось загрузить `.env`")

	}

	if env.Unmarshal("", c.DB) != nil {
		return errors.New("ошибка unmarshal процесса БД-конфига")
	}
	return nil
}
