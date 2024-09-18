# kctxsync

<p align="center">
<img src="https://img.shields.io/github/downloads/vvrnv/kctxsync/total" alt="Total Downloads">
<img src="https://img.shields.io/github/go-mod/go-version/vvrnv/kctxsync" alt="Go Version">
<a href="https://pkg.go.dev/github.com/vvrnv/kctxsync"><img src="https://pkg.go.dev/badge/github.com/vvrnv/kctxsync.svg" alt="Go Version"></a>
</p>

`kctxsync` is a command-line tool designed to synchronize the certificate and key data from a remote Kubernetes cluster's kubeconfig file to your local kubeconfig. It is particularly useful for updating the local kubeconfig when certificates or keys have changed in the remote cluster, ensuring that your local environment is always in sync with the remote server.

## Features

- **Sync kubeconfig**: Automatically fetches the kubeconfig file from a remote server via SSH and updates the local kubeconfig with any changes in the certificate-authority-data, client-certificate-data, or client-key-data.
- **Effortless updates**: No need to manually inspect or copy certificates and keys between your local and remote environments. The tool ensures that your local kubeconfig always contains the latest certificates from the remote cluster.
- **Flexible configuration**: Specify the path to your local kubeconfig file and the SSH user for remote access, allowing for a customizable workflow.

## How it works

1. The tool connects to the remote server using SSH.
2. It retrieves the remote kubeconfig from the default path (`~/.kube/config`).
3. The tool compares the certificates and keys in the remote kubeconfig with those in your local kubeconfig.
4. If differences are found, the local kubeconfig is updated with the remote data, ensuring that your environment has the correct and up-to-date credentials.

## Usage

To sync the local kubeconfig with the remote server's kubeconfig, use the following command:

```bash
kctxsync sync <context_name>
```

- <context_name>: (Optional) The name of the Kubernetes context you wish to sync.
- Optionally, you can specify the path to the local kubeconfig and the SSH user:

```bash
kctxsync sync <context_name> --config /path/to/kubeconfig --user <ssh_user>
```

If you do not provide a context and there are multiple contexts in the kubeconfig, an error will prompt you to select a context explicitly.

## Example

```bash
kctxsync sync dev --user ubuntu
```

This will connect to the remote server associated with the dev context using the ubuntu user via SSH, fetch the kubeconfig file from the server, and update your local kubeconfig with the latest certificates and keys.

## Installation

### Homebrew

```sh
brew install vvrnv/tap/kctxsync
```

### Go

```sh
go install github.com/vvrnv/kctxsync@latest
```

### Download binary

[release page link](https://github.com/vvrnv/kctxsync/releases)
