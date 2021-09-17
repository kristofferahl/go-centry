# Documentation

- [Installation](../README.md#install)
- [Getting started](../README.md#getting-started)
- [Commands](#commands)
  - [Root commands](#root-commands)
  - [Sub commands](#sub-commands)
  - [Command properties](#command-properties)
  - [Command annotations](#command-annotations)
- [Options](#options-flags)
  - [Accessing option values](#accessing-option-values)
  - [Global options](#global-options)
  - [Command options](#command-options)
  - [Option types](#option-types)
  - [Option properties](#option-properties)
  - [Option annotations](#option-annotations)
- [Arguments](#arguments)
- [Scripts](#scripts)
- [Configuration](#configuration)
  - [Metadata](#cli-metadata)
  - [Logging](#logging)
  - [Advanced](#advanced-config)
- Internal commands
- Help
- Autocompletion

## Commands

### Root commands

Root level commands are defined in the `commands` section of the manifest file (`centry.yaml`).
To define a command, two `properties` are required. The `name` is what you will be calling when using your CLI and the `path` points to the script file where the command function lives.

Here's how you would define a root level command called `get`:

_`// file: centry.yaml`_

```yaml
commands:
  - name: get
    path: ./get.sh
```

In the script file, create a function matching the `name` property.

_`// file: get.sh`_

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

_`// file: centry.yaml`_

```yaml
commands:
  - name: get
    path: ./get.sh
```

_`// file: get.sh`_

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

Adding annotations for sub commands works in the same way as for root level commands. Here's an example adding a description for the two commands created above. _Note that the full function name must be used in the annotation._

_`// file: get.sh`_

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
| ----------- | ---------------------------------------------------- | ------------- | ------- | -------- |
| Name        | The name of the command                              | `name`        | string  | true     |
| Path        | Relative path to the script containing the command   | `path`        | string  | true     |
| Description | Description of the command, displayed in help output | `description` | string  | false    |
| Help        | Usage example for the command                        | `help`        | string  | false    |
| Hidden      | When true, hides the command from help output        | `hidden`      | boolean | false    |

### Command annotations

Command annotations are used to associate metadata with a command. Annotations are defined using regular comments in bash (_a line starting with `#`_). They may be placed anywhere inside the script file and in any order you want. It is however recommended that you keep it close to your functions to act as documentation when changing your commands.

| Property    | Format                                        |
| ----------- | --------------------------------------------- |
| Description | `# centry.cmd[<command>]/description=<value>` |
| Help        | `# centry.cmd[<command>]/help=<value>`        |
| Hidden      | `# centry.cmd[<command>]/hidden=<value>`      |

## Options (flags)

Options (aka flags) are used to pass named arguments to commands. When used, `centry` will export a variable for you with the value of the option set.

### Accessing option values

Option values are made available to your commands as environment variables. Given an option named `filter`, centry sets the environment variable `FILTER` to the value provided by the option or to it's default value. The environment variable name that is used for an option can be changed by setting the `EnvName` property (see Option properties).

### Global options

Global options are made available for all commands. They are often used to to provide context for the commands you are executing. Global options are defined in the `options` section of the manifest file (`centry.yaml`). To define a global option, two properties are required. The name of the option and it's type. In general you should only specify a global option if it makes sense in the context of all commands provided by your cli.

Here's how you would define the global option `--verbose`:

_`// file: centry.yaml`_

```yaml
options:
  - name: verbose
    type: bool
    description: Use verbose logging
```

Using a global option requires you to specify the option before the name of any command. This is by design and helps drive home the fact that global options are available for all commands.

```bash
mycli --verbose command1
mycli --verbose command2
```

### Command options

Command options are, as the name suggests, scoped to commands. Therefore there is no way to define these in the manifest file. Instead you will be using `annotations` to define command options. For a full list of available annotations, see Option annotations.

Here's an example defining a `filter` option for the `get files` command:

_`// file: get.sh`_

```bash
#!/usr/bin/env bash

# centry.cmd[get:files].option[filter]/description=List only files matching the specified filter
get:files() {
  echo "listing files in the current directory (filter=${FILTER:-})..."
  if [[ "${FILTER:-}" != "" ]]; then
    ls . | grep "${FILTER:?}"
  else
    ls .
  fi
}
```

### Option types

Options have a `type` property that defines it's behavior and possible values. The currently supported option types are:

#### String option

String options are the most common type to use. It has a `name` and it's `type` set to `string`. In addition to the required properties, it is quite common to use the `default` property to set a default value for the option. See Option properties for the full list of available properties.

**Example**

_`// file: centry.yaml`_

```yaml
options:
  - name: filter
    type: string
    description: Filters output of the command
```

**Usage**: `--<option_name> <value>` or `--<option_name>=<value>`

#### Bool option

Boolean options can be used to provide a switch for behaviors in a command. As an example it could be used to turning debug logging on or off. A bool option have a value of `false` by default (this can be changed but it is not recommended). Using the default value of `false`, providing the option to your cli will tell centry to toggle that value to `true`.

**Example**

_`// file: centry.yaml`_

```yaml
options:
  - name: verbose
    type: bool
    description: Turn on verbose logging
```

_`// file: get.sh`_

```bash
#!/usr/bin/env bash

get:data() {
  echo "getting data..."
  if [[ ${VERBOSE} ]]; then
    curl --verbose http://google.com
  else
    curl http://google.com
  fi
}
```

#### Integer option

Integer options can be used to pass numbers to your commands. Things like `--max-retries=5` and `--cluster-size=3` are great examples where you might want to use an integer option. Integer options have a default value of `0` but may be set to any integer value. Passing an integer option will override the default value to the value provided.

**Example**

_`// file: get.sh`_

```bash
#!/usr/bin/env bash

# centry.cmd[get:url].option[url]/required=true
# centry.cmd[get:url].option[max-retries]/type=integer
# centry.cmd[get:url].option[max-retries]/default=3
get:url() {
  echo "Calling ${URL:?} a maximum of ${MAX_RETRIES:?} time(s)"
  echo

  local success=false
  local attempts=0
  until [[ ${success:?} == true ]]; do
    ((attempts++))

    if ! curl "${URL:?}"; then
      if [[ ${attempts:?} -ge ${MAX_RETRIES:?} ]]; then
        echo
        echo "Max retries reached... exiting!"
        return 1
      fi
      sleep 1
    else
      success=true
    fi
  done
}
```

**Usage**: `--<option_name>` or `--<option_name>=<value>`

#### Select option

Select options are a bit different. It is commonly used to have the user select one value from an array of predefined values. The user selects a value by using the matching option.

Let's dive into an example where we want the user to be able to select one of three AWS regions (eu-central-1, eu-west-1 and us-east-1). Here's how we would define that in our manifest.

_`// file: centry.yaml`_

```yaml
options:
  - name: eu-central-1
    type: select
    env_name: AWS_REGION
    description: Use eu-central-1 AWS region
  - name: eu-west-1
    type: select
    env_name: AWS_REGION
    description: Use eu-west-1 AWS region
  - name: us-east-1
    type: select
    env_name: AWS_REGION
    description: Use us-east-1 AWS region
```

On it's own, a select option provides no real value. The magic happens when we override the environment variable name that will have it's value set when a select option is provided. This essentially creates an array of valid values scoped to the specified environment variable name.

_`// file: get.sh`_

```bash
#!/usr/bin/env bash

get:lambdas() {
  : [ ${AWS_REGION:?'An AWS region must be selected using one of the predefined options'} ]
  echo "listing AWS lambda functions in region ${AWS_REGION:?}..."
  aws lambda list-functions --region ${AWS_REGION:?}
}
```

With the above command defined we can now run the following to list lambdas in the region us-east-1:

```bash
mycli get lambdas --us-east-1
```

**NOTE**:

- As no default value can be specified for select options, it's name is instead used as it's value.
- If multiple select options with the same environment variable name is specified, the last one wins.

### Option properties

| Property    | Description                                         | YAML          | Type                            | Required |
| ----------- | --------------------------------------------------- | ------------- | ------------------------------- | -------- |
| Type        | Type of option                                      | `type`        | OptionType (string/bool/select) | true     |
| Name        | Name of the option                                  | `name`        | string                          | true     |
| Short       | Short name of the option                            | `short`       | string                          | false    |
| EnvName     | Name of environment variable set for the option     | `env_name`    | string                          | false    |
| Default     | Default value of the option                         | `default`     | string                          | false    |
| Description | Description of the option, displayed in help output | `description` | string                          | false    |
| Hidden      | When true, hides the option from help output        | `hidden`      | boolean                         | false    |
| Required    | When true, marks the option as required             | `required`    | boolean                         | false    |

### Option annotations

Option annotations are used to define options for a command. Annotations are defined using regular comments in bash (a line starting with #). They may be placed anywhere inside the script file and in any order you want. It is however recommended that you keep it close to your functions to double as documentation for the command/option.

| Property    | Format                                                         |
| ----------- | -------------------------------------------------------------- |
| Type        | `# centry.cmd[<command>].option[<option>]/type=<value>`        |
| Short       | `# centry.cmd[<command>].option[<option>]/short=<value>`       |
| EnvName     | `# centry.cmd[<command>].option[<option>]/envName=<value>`     |
| Default     | `# centry.cmd[<command>].option[<option>]/default=<value>`     |
| Description | `# centry.cmd[<command>].option[<option>]/description=<value>` |
| Hidden      | `# centry.cmd[<command>].option[<option>]/hidden=<value>`      |
| Required    | `# centry.cmd[<command>].option[<option>]/required=<value>`    |

## Arguments

Anything after the last specified command or option will be passed to your command as arguments.

```bash
mycli mycommand --myoption=foo bar baz
```

Assuming `mycommand` have an option defined called `myoption`, in the example above, `bar` and `baz` would be passed as arguments. The same is true when `myoption` is left out.

### Passing flags as arguments

In some cases it is useful for flags to be passed on as arguments to a command.
In the following command we have wrapped `curl` but want to allow the use of the verbose flag, even though that is not the default behavior.

_`// file: get.sh`_

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

_`// file: centry.yaml`_

```yaml
scripts:
  - /usr/share/bash-commons-1.0.2/modules/bash-commons/src/log.sh
  - scripts/helpers.sh
  - scripts/init.sh
```

If you need something to run at the point of sourcing it, simply make it self executing like the init script below.

_`// file: scripts/init.sh`_

```bash
#!/usr/bin/env bash

init() {
  echo 'initializing the environment'
}

init "$@"
```

**NOTE**: It is important to know that naming conflicts may occur. If multiple scripts are sourced, containing functions with the same name, only the last one would be available for commands to use.

## Configuration

The `config` section of the manifest file allows you to override default values as well as describing your CLI.

### CLI metadata

The most common place to start is with the metadata properties that are used to describe your CLI to it's users. They are defined at the root of the `config` section as shown below:

```yaml
config:
  name: mycli
  description: does whatever I wan't it to
  version: 1.0.0
```

| Property    | Description            | Type   | Default                              | Required |
| ----------- | ---------------------- | ------ | ------------------------------------ | -------- |
| Name        | Name of the CLI        | string | -                                    | true     |
| Description | Description of the CLI | string | A declarative cli built using centry | false    |
| Version     | Version of the CLI     | string | -                                    | false    |

### Logging

In the `config.log` section you may change things related to logging in centry. You may want to turn on debug logging or add a prefix to the log messages printed by centry.

```yaml
config:
  name: mycli
  log:
    level: debug
    prefix: "[centry] "
```

| Property | Description                               | Type                                   | Default | Required |
| -------- | ----------------------------------------- | -------------------------------------- | ------- | -------- |
| Level    | Log level used by centry                  | LogLevel (debug/info/warn/error/panic) | info    | true     |
| Prefix   | Prefix applied to all centry log messages | string                                 | -       | false    |

### Advanced config

Also defined in the `config` section are some properties that allow even more control of how centry works.

```yaml
config:
  name: mycli
  environmentPrefix: MY_CLI_
  hideInternalCommands: false
  hideInternalOptions: false
```

| Property             | Description                                                | Type    | Default | Required |
| -------------------- | ---------------------------------------------------------- | ------- | ------- | -------- |
| EnvironmentPrefix    | Prefix used when exporting environment variables in centry | string  | -       | false    |
| HideInternalCommands | Hides internal centry commands from help output            | boolean | true    | false    |
| HideInternalOptions  | Hides internal centry options from help output             | boolean | true    | false    |
