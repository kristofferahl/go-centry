# centry

Declarative **command line builder** for teams and one man bands

## Use cases
- Build a feature rich CLI from scratch using only scripts and yaml
- Bundle existing scripts and tools as a single CLI
- Encode best practices and team conventions for using scripts and dev/ops tools

## Feature highlights
- **Declarative**: Add new commands and options and run them at once. No need to re-compile.
- **Unified syntax**: Provides a standard for using commands and options.
- **Supports multi level commands**: `mycli get status` and `mycli get files`
- **Supports multi level options (flags)**: `mycli --dev get status --out json`
- **Contextual help**: `mycli get --help`
- **Autocomplete**: Bash-completions of commands and options
- **Highly configurable**: Sensible defaults, lots of choises
- **Easy setup**: Download centry, create a manifest file and you are good to go


## Install

### Mac
```bash
curl -L https://github.com/kristofferahl/go-centry/releases/download/v1.0.0-prerelease2/go-centry_1.0.0-prerelease2_Darwin_x86_64.tar.gz | tar -xzv -C /usr/local/bin/
```

### Linux
```bash
curl -L https://github.com/kristofferahl/go-centry/releases/download/v1.0.0-prerelease2/go-centry_1.0.0-prerelease2_Linux_x86_64.tar.gz | tar -xzv -C /usr/local/bin/
```

## Getting started

**The documentation and examples below assumes that**
1. You are running `bash` version 3.2 or later
1. You have "installed" the `go-centry_*` binary for your OS and made it available in your path as `mycli` (by renaming the file)
1. You have created an empty directory to hold your commands and manifest file

## Setup
1. Create the manifest file for the CLI and name it `centry.yaml` by running the following command in your shell.
    ```bash
    echo "commands: []
    config:
      name: mycli" > centry.yaml
    ```
2. Verify that it's working by running
    ```
    mycli --help
    ```
   This should display the contextual help for the cli and the name **mycli** at the top.

## The manifest file
This is where you define root level commands and options, do configuration overrides and import scripts to be available for all your commands.

By default, `centry` will look for a `centry.yaml` file in the **current directory**. You may change the location and name of the manifest file but this requires you to let centry know where to find it. This can be done by setting the environment variable `CENTRY_FILE` or by way of passing `--centry-file <path>` as the **first** argument.

## Commands
In `centry`, commands are simple shell scripts with a matching function name in it.

Let's start by creating a file called `hello.sh` with the following content.

*`// file: hello.sh`*
```
#!/usr/bin/env bash

hello() {
  echo 'Hello centry'
}
```

Before you can use the `hello` function as a command, you need to tell `centry` where to find it. Open `centry.yaml` in an editor of choise and modify it to look like this:

*`// file: centry.yaml`*
```yaml
commands:
  - name: hello
    path: ./hello.sh
    description: Says hello

config:
  name: mycli
```

You should now be able to able to run the command.
```bash
mycli hello ↵
Hello centry
```

## Options
In `centry`, options are flags you use to pass named arguments to your command functions. This enables easier discovery of your cli and less friction for users. They may be specified in long (`--option`) or short (`-o`) form.

Let's add a `--name` option to the hello command. This is done by adding `annotations` in your script. Edit `hello.sh` to look like this.

*`// file: hello.sh`*
```
#!/usr/bin/env bash

# centry.cmd[hello].option[name]/type=string
hello() {
  echo "Hello ${NAME}"
}
```

Running the `hello` command again would look like this
```bash
mycli hello ↵
Hello
```

To pass a name to be echoed back to you, call the command with the `--name` option.
```bash
mycli hello --name William ↵
Hello William
```

If you want to add a description for the `--name` option you should add an additional annotation to the `hello.sh` file.

*`// file: hello.sh`*
```
#!/usr/bin/env bash

# centry.cmd[hello].option[name]/type=string
# centry.cmd[hello].option[name]/description=Name to be greeted with
hello() {
  echo "Hello ${NAME}"
}
```
Displaying the contextual help (using the `--help` option) should now look something like this.
```bash
mycli hello --help ↵
NAME:
   mycli hello - Says hello

USAGE:
   mycli hello [command options] [arguments...]

OPTIONS:
   --name value  Name to be greeted with
   --help, -h    Show help (default: false)
```

## Arguments
A command may also accept any number of arguments. All arguments not matching an option of a command will be passed on to the function.

*`// file: hello.sh`*
```
#!/usr/bin/env bash

# centry.cmd[hello].option[name]/type=string
# centry.cmd[hello].option[name]/description=Name to be greeted with
hello() {
  echo "Hello ${NAME}"
  echo "Arguments (${#*}): ${*}"
}
```

NOTE: Arguments must always be passed after the last option.
```bash
mycli hello --name William arg1 arg2 ↵
Hello William
Arguments (2): arg1 arg2
```

## Autocomplete

**NOTE: Only available for Bash**

To make discovery of `mycli` easier, we may want to enable bash completions. Follow the steps below to set it up.
```bash
curl -o bash_autocomplete https://raw.githubusercontent.com/kristofferahl/go-centry/master/bash_autocomplete
PROG=mycli source bash_autocomplete
```

Now, let try it out by typing `mycli` followed by a space and then hit `tab`. This will display any command available at the root level. If there is only one, the command name will be autocompleted. It works for options too.
```bash
mycli -- ➡
--centry-config-log-level  --centry-quiet             --help
```

