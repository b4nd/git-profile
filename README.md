# Git Profile

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go](https://github.com/b4nd/git-profile/actions/workflows/build.yml/badge.svg)](https://github.com/b4nd/git-profile/actions/workflows/build.yml)

## Overview
Git Profile is a command-line application developed in Go, designed to manage multiple Git profiles efficiently. It provides a suite of commands to create, update, delete, and switch between different user profiles, making it easier to handle various identities across different Git repositories. This tool is particularly useful for developers who work on multiple projects with different user credentials.

## Commands Documentation

| Command | Description |
|---------|-------------|
| `git-profile current` | Displays the currently active profile. |
| `git-profile [delete\|del]` | Deletes a specified profile from the system. |
| `git-profile get` | Retrieves details of a specific profile. |
| `git-profile [list\|ls]` | Lists all available profiles. |
| `git-profile set` | Sets or updates a profile configuration. |
| `git-profile use` | Switches to a specific profile for operations. |
| `git-profile amend` | Updates email and name of the current profile last commit. |
| `git-profile version` | Displays the current version of the application. |
| `git-profile help` | Displays help information for the application. |
| `git-profile completion` | Generates shell completion scripts. |

## Installation

### Manual installation:

1. Download the latest release from the [releases page](https://github.com/b4nd/git-profile/releases)
2. Extract the archive and move the binary to a directory in your PATH.

### Linux:

```bash
curl -sL https://github.com/b4nd/git-profile/releases/download/v0.1.0/git-profile-v0.1.0-linux-amd64 -o git-profile
chmod +x git-profile 
sudo mv git-profile /usr/local/bin/
```

## Usage

To use the application, run the following command:

```bash
git profile [command] [flags]
```

## Environment variables

| Variable | Description |
|---------|-------------|
| `GIT_PROFILE_PATH` | The path to the directory where the profiles are stored. Default is `$HOME/.git-profile`. |

## Requirements

- [Go](https://golang.org/)
- [Taskfile](https://taskfile.dev/)
- [Git](https://git-scm.com/)

## Build

1. Clone the repository:
   ```bash
   git clone https://github.com/b4nd/git-profile.git
   ```
2. Change into the project directory:
   ```bash
   cd git-profile
   ```
3. Install the dependencies:
   ```bash
   go mod tidy
   ```
4. Run the tests:
   ```bash
   task test
   ```
5. Run the following command to install the application:
   ```bash
   task build
   ```
build output will be in the `bin` directory.

## Example usage during development

Below are examples of how to use each command:

```bash
task run -- set \
   --workspace company \
   --name "Your Name" \
   --email "name@example.com"
task run -- list
task run -- use company
task run -- current
```
## Contributing

1. Fork the repository.
2. Create a new branch:
   ```bash
   git checkout -b feature-name
   ```
3. Commit your changes:
   ```bash
   git commit -m 'Add new feature'
   ```
4. Push to the branch:
   ```bash
   git push origin feature-name
   ```
5. Open a pull request.

## Authors

- [Juan Manuel Garcia](https://github.com/b4nd/me)

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

