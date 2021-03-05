# Documentation

- Getting started (see README)
- Commands
    - Root commands
    - Sub commands
    - Command properties
    - Command annotations
- Options
    - Global options
    - Command options
    - Option properties
    - Option annotations
- Arguments
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

### Root commands

Root level commands are defined in the `commands` section of the manifest file (`centry.yaml`).
To define a command, two `properties` are required. The `name` is what you will be calling when using your CLI and the `path` points to the script file where the command function lives.

Here's how you would define a root level command called `get`:

*`// file: centry.yaml`*
```yaml
commands:
    - name: get
      path: ./get.sh
```

In the script file, create a function matching the `name` property.

*`// file: get.sh`*
```bash
#!/usr/bin/env bash

get() {
  echo "getting stuff for the fun of it"
}
```

There are additional properties that may be set for a command (see `Command properties`).
You may also choose to specify some of the properties using annotations (see `Command annotations`).
What strategy you choose is entirely up to you but the root level commands must always be partially specified in the manifest file.

### Sub commands

Sub commands are exclusively defined in scripts. Creating a sub command is as easy as including the special character colon (`:`) in a script function name. Let's say you have already defined a root level command called `get` but wanted to define two commands that have `get` as their parent. Simply create two functions named `get:` and suffix it with the desired name of the sub command.

*`// file: centry.yaml`*
```yaml
commands:
    - name: get
      path: ./get.sh
```

*`// file: get.sh`*
```bash
#!/usr/bin/env bash

get:data() {
  echo "getting the latest data..."
}

get:time() {
  echo "the time is $(date +"%T")"
}
```

The script above defins two subcommands, `data` and `time`. The can now be executed by calling your CLI like below.

```bash
mycli get data
mycli get time
```

Adding annotations for sub commands works in the same way as for root level commands. Here's an example adding a description for the two commands created above. *Note that the full function name must be used in the annotation.*

```bash
#!/usr/bin/env bash

# centry.cmd[get:data]/description=Get's you data
get:data() {
  echo "getting the latest data..."
}

# centry.cmd[get:time]/description=Displays the current time
get:time() {
  echo "the time is $(date +"%T")"
}
```

### Command properties

| Property    | Description                                          | YAML key      | Type    | Required |
|-------------|------------------------------------------------------|---------------|---------|----------|
| Name        | The name of the command                              | `name`        | string  | true     |
| Path        | Relative path to the script containing the command   | `path`        | string  | true     |
| Description | Description of the command, displayed in help output | `description` | string  | false    |
| Help        | Usage example for the command                        | `help`        | string  | false    |
| Hidden      | When true, hides the command from help output        | `hidden`      | boolean | false    |


### Command annotations

Command annotations are used to associate metadata with a command. Annotations are defined using regular comments in bash (*a line starting with `#`*). They may be placed anywhere inside the script file and in any order you want. It is however recommended that you keep it close to your functions to act as documentation when changing your commands.

| Property    | Format                                        |
|-------------|-----------------------------------------------|
| Description | `# centry.cmd[<command>]/description=<value>` |
| Help        | `# centry.cmd[<command>]/help=<value>`        |
| Hidden      | `# centry.cmd[<command>]/hidden=<value>`      |

## Options

### Option properties

| Property    | Description                                         | YAML          | Type                            | Required |
|-------------|-----------------------------------------------------|---------------|---------------------------------|----------|
| Type        | Type of option                                      | `type`        | OptionType (string/bool/select) | true     |
| Name        | Name of the option                                  | `name`        | string                          | true     |
| Short       | Short name of the option                            | `short`       | string                          | false    |
| EnvName     | Name of environment variable set for the option     | `env_name`    | string                          | false    |
| Default     | Default value of the option                         | `default`     | string                          | false    |
| Description | Description of the option, displayed in help output | `description` | string                          | false    |
| Hidden      | When true, hides the option from help output        | `hidden`      | boolean                         | false    |

### Option annotations

| Property    | Format                                                         |
|-------------|----------------------------------------------------------------|
| Type        | `# centry.cmd[<command>].option[<option>]/type=<value>`        |
| Short       | `# centry.cmd[<command>].option[<option>]/short=<value>`       |
| EnvName     | `# centry.cmd[<command>].option[<option>]/envName=<value>`     |
| Default     | `# centry.cmd[<command>].option[<option>]/default=<value>`     |
| Description | `# centry.cmd[<command>].option[<option>]/description=<value>` |
| Hidden      | `# centry.cmd[<command>].option[<option>]/hidden=<value>`      |

## Arguments

Anything after the last specified command or option will be passed to your command as arguments.

```bash
mycli mycommand --myoption=foo bar baz
```
Assuming `mycommand` have an option defined called `myoption`, in the example above, `bar` and `baz` would be passed as arguments. The same is true when `myoption` is left out.

### Passing flags as arguments
In some cases it is useful for flags to be passed on as arguments to a command.
In the following command we have wrapped `curl` but want to allow the use of the verbose flag, even though that is not the default behavior.

*`// file: get.sh`*
```bash
#!/usr/bin/env bash

# centry.cmd[get:data]/description=Get's data from a URL
# centry.cmd[get:data].option[url]/description=The URL to get data from
get:data() {
  echo "getting the data from ${URL:?}..."
  curl "${URL:?}" "$@"
}
```

The command below will fail since a verbose option is not defined.

```bash
mycli get data --url http://google.com --verbose
```

 To achive the desired behaviour we need to tell centry to **stop processing arguments**. This can be done by adding a `--` when calling the command.

```bash
mycli get data --url http://google.com -- --verbose
```

## Scripts
Before executing a command, centry can import helper functions and run common setup tasks for the environment the command executes in. This is done by specifying an array of file paths in the `scripts` section that centry will [source](https://linuxize.com/post/bash-source-command/), in the specified order. This makes sharing functions across commands easier and more predictable while keeping things DRY.

Comman use-cases include:
- Sourcing of functions from script libraries and wrapper scripts
- Installing missing dependencies and downloading files
- Setting environment variables
- Authentication and authorization

Here's an example:

*`// file: centry.yaml`*
```yaml
scripts:
  - /usr/share/bash-commons-1.0.2/modules/bash-commons/src/log.sh
  - scripts/helpers.sh
  - scripts/init.sh
```

If you need something to run at the point of sourcing it, simply make it self executing like the init script below.

*`// file: scripts/init.sh`*
```bash
#!/usr/bin/env bash

init() {
  echo 'initializing the environment'
}

init "$@"
```

**NOTE**: It is important to know that naming conflicts may occur. If multiple scripts are sourced, containing functions with the same name, only the last one would be available for commands to use.