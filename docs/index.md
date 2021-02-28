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

| Property        | Description                                          | YAML key      | Type   | Required |
|-------------|------------------------------------------------------|---------------|--------|----------|
| Name  | The name of the command   | `name`  | string  | true  |
| Path   | Relative path to the script containing the command   | `path`  | string  | true  |
| Description | Description of the command, displayed in help output | `description` | string | false    |
| Help   | Usage example for the command  | `help`  | string  | false  |
| Hidden   | When set to true, hides the command from help output  | `hidden`  | boolean   | false  |


### Annotations

| Property    | Format                                        |
|-------------|-----------------------------------------------|
| Description | `# centry.cmd[<command>]/description=<value>` |
| Help        | `# centry.cmd[<command>]/help=<value>`        |
| Hidden      | `# centry.cmd[<command>]/hidden=<value>`      |
