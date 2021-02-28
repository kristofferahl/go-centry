# Documentation

- Getting started (see README)
- Commands
    - Properties
    - Annotations
    - Root level commands
    - Sub commands
- Options
    - Properties
    - Annotations
    - Global options
    - Command options
- Scripts
- Configuration
    - Application properties
        - Name
        - Description
        - Version
    - Log
        - Level
        - Prefix
- Help
- Autocompletion

## Commands

### Properties

| Property    | Description                                          | YAML key      | Type    | Required |
|-------------|------------------------------------------------------|---------------|---------|----------|
| Name        | The name of the command                              | `name`        | string  | true     |
| Path        | Relative path to the script containing the command   | `path`        | string  | true     |
| Description | Description of the command, displayed in help output | `description` | string  | false    |
| Help        | Usage example for the command                        | `help`        | string  | false    |
| Hidden      | When true, hides the command from help output        | `hidden`      | boolean | false    |


### Annotations

| Property    | Format                                        |
|-------------|-----------------------------------------------|
| Description | `# centry.cmd[<command>]/description=<value>` |
| Help        | `# centry.cmd[<command>]/help=<value>`        |
| Hidden      | `# centry.cmd[<command>]/hidden=<value>`      |

## Options

### Properties

| Property    | Description                                         | YAML key      | Type                            | Required |
|-------------|-----------------------------------------------------|---------------|---------------------------------|----------|
| Type        | Type of option                                      | `type`        | OptionType (string/bool/select) | true     |
| Name        | Name of the option                                  | `name`        | string                          | true     |
| Short       | Short name of the option                            | `short`       | string                          | false    |
| EnvName     | Name of environment variable set for the option     | `env_name`    | string                          | false    |
| Default     | Default value of the option                         | `default`     | string                          | false    |
| Description | Description of the option, displayed in help output | `description` | string                          | false    |
| Hidden      | When true, hides the option from help output        | `hidden`      | boolean                         | false    |

### Annotations

| Property    | Format                                                         |
|-------------|----------------------------------------------------------------|
| Type        | `# centry.cmd[<command>].option[<option>]/type=<value>`        |
| Short       | `# centry.cmd[<command>].option[<option>]/short=<value>`       |
| EnvName     | `# centry.cmd[<command>].option[<option>]/envName=<value>`     |
| Default     | `# centry.cmd[<command>].option[<option>]/default=<value>`     |
| Description | `# centry.cmd[<command>].option[<option>]/description=<value>` |
| Hidden      | `# centry.cmd[<command>].option[<option>]/hidden=<value>`      |
