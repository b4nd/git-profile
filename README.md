# Git Profile

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go](https://github.com/b4nd/git-profile/actions/workflows/build.yml/badge.svg)](https://github.com/b4nd/git-profile/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/b4nd/git-profile)](https://goreportcard.com/report/github.com/b4nd/git-profile)
[![codecov](https://codecov.io/gh/b4nd/git-profile/graph/badge.svg?token=LR1HAJ71CY)](https://codecov.io/gh/b4nd/git-profile)
[![Release](https://img.shields.io/github/release/b4nd/git-profile.svg)](https://github.com/b4nd/git-profile/releases/latest)
[![GitHub Releases Stats of git-profile](https://img.shields.io/github/downloads/b4nd/git-profile/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=b4nd&repository=git-profile)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/3724fd18d1934563857274672e02f3fa)](https://app.codacy.com/gh/b4nd/git-profile/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/3724fd18d1934563857274672e02f3fa)](https://app.codacy.com/gh/b4nd/git-profile/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

#### SonarCloud 

[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=b4nd_git-profile&metric=bugs)](https://sonarcloud.io/summary/new_code?id=b4nd_git-profile)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=b4nd_git-profile&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=b4nd_git-profile)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=b4nd_git-profile&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=b4nd_git-profile)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=b4nd_git-profile&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=b4nd_git-profile)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=b4nd_git-profile&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=b4nd_git-profile)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=b4nd_git-profile&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=b4nd_git-profile)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=b4nd_git-profile&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=b4nd_git-profile)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=b4nd_git-profile&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=b4nd_git-profile)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=b4nd_git-profile&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=b4nd_git-profile)

## Overview

Git Profile is a command-line application developed in Go, designed to efficiently manage multiple Git profiles. It provides a suite of commands to create, update, delete, and switch between different user profiles, making it easier to handle various identities across diverse Git repositories. The commands also feature an interactive mode that prompts you for the necessary details, simplifying the setup process. This tool is particularly useful for developers who work on multiple projects with different user credentials.


### Example Use Case

Imagine a developer working on both open-source and corporate projects. They need to switch between different Git profiles seamlessly to ensure commits are associated with the correct email and username. Instead of manually changing Git configurations every time, they can use `git-profile` to quickly switch between predefined profiles, improving workflow efficiency.
Once you’ve configured a project profile, git-profile will remember it the next time you work on that project, saving you the hassle of reconfiguration.

#### Profile Storage

By default, git-profile stores profiles in `$HOME/.gitprofile`. However, you can also store them locally by using the `--local` flag, which places a `.gitprofile` file in the current folder. This feature is especially useful for keeping project-specific settings right inside the repository.

#### Where the selected profile is stored

git-profile uses the .git/config file in each repository to store the selected profile. This way, there’s no need to reconfigure the profile every time you work on that repository, and it also ensures that the local name and user remain consistent for each project.

![git-profile](https://raw.githubusercontent.com/b4nd/git-profile/main/doc/git-profile.gif)

## Commands Documentation

| Command                     | Description                                                |
| --------------------------- | ---------------------------------------------------------- |
| `git-profile current`       | Displays the currently active profile.                     |
| `git-profile [delete\|del]` | Deletes a specified profile from the system.               |
| `git-profile get`           | Retrieves details of a specific profile.                   |
| `git-profile [list\|ls]`    | Lists all available profiles.                              |
| `git-profile set`           | Sets or updates a profile configuration.                   |
| `git-profile use`           | Switches to a specific profile for operations.             |
| `git-profile amend`         | Updates email and name of the current profile last commit. |
| `git-profile version`       | Displays the current version of the application.           |
| `git-profile help`          | Displays help information for the application.             |
| `git-profile completion`    | Generates shell completion scripts.                        |

## Installation

### Linux

```bash
curl -sL https://github.com/b4nd/git-profile/releases/download/v0.1.3/git-profile-v0.1.3-linux-amd64 -o git-profile
chmod +x git-profile 
mv git-profile /usr/local/bin/
```

### macOS

```bash
curl -sL https://github.com/b4nd/git-profile/releases/download/v0.1.3/git-profile-v0.1.3-darwin-amd64 -o git-profile
chmod +x git-profile 
mv git-profile /usr/local/bin/
```

### Windows

1. Download the latest Windows executable from the [releases page](https://github.com/b4nd/git-profile/releases).
2. Extract the archive.
3. Move the `git-profile-v0.1.3-darwin-amd64.exe` file to a directory in your system `PATH` and rename it to `git-profile.exe`.
4. Optionally, add the directory to the system `PATH` environment variable for easier access.

```powershell
[System.Environment]::SetEnvironmentVariable("Path", $Env:Path + ";C:\\path\\to\\git-profile", [System.EnvironmentVariableTarget]::User)
```

## Usage

To use the application, run the following command:

```bash
git profile [command] [flags]
```

### More Examples

- **Create a new profile for a personal project:**

  ```bash
  git profile set --workspace personal --name "Your Name" --email "name@example.com"
  ```

  This command sets up a new profile named `personal` with the given credentials.

- **List all existing profiles:**

  ```bash
  git profile list
  ```

  Displays all available profiles currently stored.

- **Use a specific profile:**

  ```bash
  git profile use personal
  ```

  Switches to the `personal` profile, applying its Git credentials.

- **Check the currently active profile:**

  ```bash
  git profile current
  ```

  Shows which profile is currently in use.

- **Amend the last commit with the active profile's details:**

  ```bash
  git profile amend
  ```

  Updates the latest commit with the email and name of the currently active profile.

- **Remove a profile:**

  ```bash
  git profile delete personal
  ```

  Deletes the `personal` profile from the system.

## Environment variables

| Variable           | Description                                                                               |
| ------------------ | ----------------------------------------------------------------------------------------- |
| `GIT_PROFILE_PATH` | The path to the directory where the profiles are stored. Default is `$HOME/.gitprofile`. |

### Configuring GIT\_PROFILE\_PATH in `.zshrc` or `.bashrc`

If you want to specify a custom location for the Git Profile configuration, you can set the `GIT_PROFILE_PATH` environment variable in your shell configuration file.

For Windows Subsystem for Linux (WSL), you can add the following lines to your `~/.zshrc` or `~/.bashrc`:

```bash
# Example of a shared Git profile on Windows through WSL
export GIT_PROFILE_PATH="/mnt/c/Users/<USER>/.gitprofile"
```

After adding the line, apply the changes by running:

```bash
source ~/.zshrc  # If using zsh
source ~/.bashrc  # If using bash
```

## DevContainer Support

This project includes support for **DevContainers**, allowing developers to quickly set up a consistent development environment using **VS Code Remote - Containers** or **GitHub Codespaces**.

### How to Use

1. Ensure you have **Docker** installed and running.
2. Open the project in **VS Code**.
3. Install the **Dev Containers** extension if you haven't already.
4. Open the Command Palette (`Ctrl+Shift+P` or `Cmd+Shift+P` on macOS) and select `Remote-Containers: Reopen in Container`.

This will automatically set up all dependencies and configurations needed for development.

## Requirements

- [Go](https://golang.org/) v1.23.6
- [Taskfile](https://taskfile.dev/) v3.41.0
- [golangci-lint](https://golangci-lint.run/) v1.63.4
- [gosec](https://github.com/securego/gosec/) v2.22.0
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

Below are examples of how to use each command, along with explanations of their purpose:

- **Set a new profile:**

  ```bash
  task run -- set \
     --workspace company \
     --name "Your Name" \
     --email "name@example.com"
  ```

  This command creates a new Git profile under the workspace `company`, assigning the specified name and email.

- **List all profiles:**

  ```bash
  task run -- list
  ```

  Displays all available Git profiles configured in the system.

- **Switch to a specific profile:**

  ```bash
  task run -- use company
  ```

  Activates the Git profile associated with `company`, ensuring that subsequent Git commits use the corresponding credentials.

- **Check the currently active profile:**

  ```bash
  task run -- current
  ```

  Shows the details of the currently active Git profile, including name and email.

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

## Acknowledgements

[![SonarQube Cloud](https://sonarcloud.io/images/project_badges/sonarcloud-light.svg)](https://sonarcloud.io/summary/new_code?id=b4nd_git-profile)