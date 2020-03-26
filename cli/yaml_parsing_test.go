package main

import (
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestYamlParsing(t *testing.T) {
	var config = Config{}
	data, _ := ioutil.ReadFile("deploy.yaml")
	var err = yaml.Unmarshal(data, &config)
	assert.Nil(t, err)
	assert.EqualValues(t, Config{
		Version:     "v1beta1",
		Name:        "app",
		App_version: "1.0.0",
		Images: map[string]Image{
			"app": {
				Dockerfile: "dockerfiles/app.dockerfile",
			},
		},
		Services: map[string]Service{
			"app": {
				Image: "@images.app",
				Depends_on: []string{
					"database",
				},
				Labels: []map[string]string{
					{"variation": "promotion"},
				},
				Ingress: Ingress{
					TLS: true,
					Rules: []map[string]string{{
						"host": "${HOSTNAME}",
					}},
				},
				Ports: []Port{{
					Port: "8080:80",
				}},
				Environment: map[string]string{
					"AUTH_HOSTNAME": "@dependencies.auth",
					"DB_NAME":       "${DB_NAME}",
					"DB_PASSWORD":   "${DB_PASSWORD}",
					"DB_USER":       "${DB_USER}",
				},
				Deploy: Deploy{
					Strategy: "ab",
					Variations: []Variation{{
						Labels: []map[string]string{{
							"variation": "default",
						}},
						Weight: 50,
					}, {
						Labels: []map[string]string{{
							"variation": "promotion",
						}},
						Weight: 50,
					}},
				},
			},
			"database": {
				Chart: Chart{
					Name:    "stable/postgresql",
					Version: "8.1.4",
					Values: map[string]string{
						"postgresqlDatabase": "${DB_NAME}",
						"postgresqlUsername": "${DB_USER}",
						"postgresqlPassword": "${DB_PASSWORD}",
					},
				},
			},
		},
	}, config)
}
