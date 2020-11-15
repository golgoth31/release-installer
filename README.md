# release-installer

## What is it ?

Many "devops" tools are single, statically linked, binary. Package managers are not always up to date and brew is almost made for Mac users. This tools **ri** is, it-self, a single binary, easy to install and esay to use to install binaries found on github for exemple.

## How to use it

For the first run, download the needed binary file from the release page and run the init command. Once installed, ri will be able to auto update (binary and yaml definitions).

```bash
Usage:
  ri [command]

Available Commands:
  help        Help about any command
  init        Initialize release-installer
  install     Install one release
  release     List available version of a release
  releases    List available releases
  update      Update the releases definitions
  version     Show the ro=i version

Flags:
      --config string   config file (default is $HOME/.release-installer/release-installer.yaml)
      --debug           debug
  -h, --help            help for ri

Use "ri [command] --help" for more information about a command.
```

### Install a release

```txt
Usage:
  ri install [release] [flags]

Flags:
      --apiVersion string   Set this install apiVersion (default "install/v1")
  -a, --arch string         Release binary architecture (default "amd64")
  -d, --default             Set this install as default
  -h, --help                help for install
  -o, --os string           Release binary OS (default "linux")
  -p, --path string         Destination to install binary in, should be set in your "$PATH" (default "~/bin")
  -v, --version string      Release version

Global Flags:
      --config string   config file (default is $HOME/.release-installer/release-installer.yaml)
      --debug           debug
```

Exemple:

```bash
ri install terrafom -v 0.13.2 -d
```

Result:

```txt

 Installing release "terraform"

» Version: 0.13.2
» OS:      linux
» Arch:    amd64
» Default: true
» Path:    ~/bin

 Release files

» Checksum file: https://releases.hashicorp.com/terraform/0.13.2/terraform_0.13.2_SHA256SUMS
» Archive file:  https://releases.hashicorp.com/terraform/0.13.2/terraform_0.13.2_linux_amd64.zip

 Downloading files

√ [========================================================================================] terraform_0.13.2_SHA256SUMS (731.43 KB/s)
√ [===================================================================================] terraform_0.13.2_linux_amd64.zip (543.31 KB/s)

√ File saved as: /home/david/bin/terraform_0.13.2

» Creating symlink: /home/david/bin/terraform

√ Done

 √ Release installed
```
