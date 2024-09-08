package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/loupe-co/{{.repoName}}/internal/config"
	"github.com/loupe-co/{{.repoName}}/internal/handlers"
	configUtil "github.com/loupe-co/go-common/config"
	"github.com/loupe-co/go-common/env"
	"github.com/loupe-co/go-common/fixtures"
)

var testHandlers *handlers.Handlers

func setup() error {
	if err := env.LoadFromFile("./local.env"); err != nil {
		return err
	}

	if err := fixtures.InitTestFixtures("../fixtures"); err != nil {
		return err
	}

	testConfig := config.Config{}
	if err := configUtil.LoadFromENV(&testConfig); err != nil {
		return err
	}

	testHandlers = handlers.New(testConfig)

	return nil
}

func teardown() error {
	return nil
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		fmt.Println(err)
		if err := teardown(); err != nil {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	code := m.Run()
	if err := teardown(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(code)
}
