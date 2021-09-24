# Adding new flag types

## internal/pkg/cmd/

1. Create constant for new type in option.go (IntegerOption)
2. Add StringToOptionType test for new type and make it pass

## schema & test data

1. Add value type string to enum in schemas/manifest.json
1. Add option example to manifest_test_valid.yaml in test/data
1. Add option example to runtime_test.yaml in test/data

## optionsSetToFlags

1. Add to switch case for handling the type conversion to an cli.Flag in optionsSetToFlags()

## optionsSetToEnvVars

1. Add EnvironmentVariableType representing the type to environment.go
1. Add to switch case for handling the type conversion to an environment variable in optionsSetToEnvVars()

## required options

1. Add test case for required options in runtime_test.go ("invoke without required option")
1. Define option for "optiontest:required" in option_test.sh
1. Run tests and make it pass
