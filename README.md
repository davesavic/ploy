# Ploy
**Ploy** is a tool designed to automate and streamline the process of software deployment and task execution across various environments. It offers a simplified and flexible solution for managing complex deployment workflows and executing predefined tasks on local or remote servers.

### Key Features
- **Simplified Deployment**: Define and execute deployment processes with ease, across multiple servers and environments.
- **Task Automation**: Automate repetitive tasks with customizable scripts.
- **Configurable Pipelines**: Set up pipelines for different deployment or task scenarios, ensuring consistency and reliability.
- **Parameter Substitution**: Dynamically substitute parameters in tasks, making your scripts more flexible and environment-agnostic.


### Installation
Download the latest release for your os from the [releases page](https://github.com/davesavic/ploy/releases)

### Usage
`ploy [command]`

### Commands
- `init` - Initialize a new template ploy script.
- `run [options] [pipeline]...` - Run pipelines on their respective servers (provide `-l` to run them locally).
- `help` - Display the help message.

### Configuration structure
```json
{
    "params": {
        // User-defined parameters for tasks
        "message": "hello, world!"
    },
    "servers": {
        // Server configurations for remote execution
        "staging": {
            "host": "111.111.111.111",
            "port": 22,
            "user": "ploy",
            "private-key": "/home/user/.ssh/id_rsa"
        }
    },
    "tasks": {
        // Task definitions for automation
        "print-message": [
            "echo '{{message}}'"
        ]
    },
    "pipelines": {
        // Pipeline configurations for deployment or task execution
        "say-hello": {
            "servers": ["staging"],
            "tasks": ["print-message"]
        }
    }
}
```

### Auto populated parameters
- `{{timestamp}}` -  Inserts the current timestamp into tasks.
- More to come...

### Getting started
To get started with Ploy, follow these steps:

1. **Install Ploy**: Download and install the latest version from the [releases page](https://github.com/davesavic/ploy/releases).
2. **Initialize Configuration**: Run `ploy init` to generate a baseline configuration file.
3. **Customize Configuration**: Edit the generated `configuration.json` to suit your deployment and task requirements.
4. **Run Ploy**: Execute your configurations using `ploy run [pipelines (run multiple separated by a space)]`.

### Recipes
- [Zero Downtime Deployments](https://github.com/davesavic/ploy/blob/master/recipes/zero-downtime-deployment.json)

### Contributing
Contributions to Ploy are welcome! 
If you have suggestions, bug reports, or want to contribute code, please visit create the relevant report.