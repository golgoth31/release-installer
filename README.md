# release-installer

## What is it ?

Many "devops" tools are single, statically linked, binary. Package managers are not always up to date and brew is almost made for Mac users. This tools **ri** is, it-self, a single binary, easy to install and esay to use to install binaries found on github for exemple.

## Descriptions

By default, the known descriptions will be put in ~/.release-installer/releases. You can add yours, it will be taken into account. But, if a software you use is not available, don't hesitate to propose a yaml description [here](https://github.com/golgoth31/release-installer-definitions).

## How to use it

For the first run, download the needed binary file from the release page and run the init command. Once installed, **_ri_** will be able to auto update (binary and yaml definitions).

```txt
$ ri -h

A tool to download and install binaries

Usage:
  ri [command]

Available Commands:
  help        Help about any command
  init        Initialize release-installer
  install     Install one release
  list        List available releases or versions
  remove      Remove a specific release version
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
$ ri install -h

Install one release

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

√ [==================================================] terraform_0.13.2_SHA256SUMS (731.43 KB/s)
√ [=============================================] terraform_0.13.2_linux_amd64.zip (543.31 KB/s)

√ File saved as: /home/david/bin/terraform_0.13.2

» Creating symlink: /home/david/bin/terraform

√ Done

 √ Release installed
```

### list releases

```txt
$ ri list -h

Usage:
  ri list [release name] [flags]

Flags:
  -h, --help         help for list
  -i, --installed    Show installed releases only
  -n, --number int   Number of releases or versions to show (default 5)

Global Flags:
      --config string   config file (default is $HOME/.release-installer/release-installer.yaml)
      --debug           debug
```

#### list all known releases

```bash
$ ri list

 Available releases

√ dive (v0.9.2)
» eksctl
√ gitcomm (v0.3.4)
» github-cli
√ goreleaser (v0.147.2)
√ helm (v3.4.0)
√ istioctl (1.7.4)
» kubectl-argo-rollouts
» skaffold
√ stern (v1.13.1)
√ terraform (0.13.5)
```

#### list all versions for one releases

```bash
$ ri list terraform

 Available versions for release "terraform"

» v0.14.6
» v0.15.0-alpha20210127
» v0.14.5
» v0.15.0-alpha20210107
» v0.14.4
```

### Update myself

```bash
$ ri update -h

Update the releases definitions

Usage:
  ri update [flags]

Flags:
  -f, --force         Force update
  -h, --help          help for update
  -p, --path string   Destination to install binary in, should be set in your "$PATH" (default "~/bin")

Global Flags:
      --config string   config file (default is $HOME/.release-installer/release-installer.yaml)
      --debug           debug
```
