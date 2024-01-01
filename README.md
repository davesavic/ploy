# Ploy
### A simplified deployment and task automation tool.

**Usage**: `ploy [command]`

**Commands**:
- `init` - Initialize a new template ploy script.
- `run [options] [pipeline]...` - Run pipelines on their respective servers (provide -l to run them locally).
- `help` - Display the help message.

**Configuration structure**:
```json
{
    "params": {
        // Parameters to be populated within tasks
        // Keys are the parameter names, values are the parameter values
        "message": "hello, world!"
    },
    "servers": {
        // Servers to be used in pipelines
        // Keys are the server names, values are the server configurations
        "staging": {
            "host": "111.111.111.111",
            "port": 22,
            "user": "ploy",
            "private-key": "/home/user/.ssh/id_rsa"
        }
    },
    "tasks": {
        // Tasks to be used in pipelines
        // Keys are the task names, values are the task commands
        "print-message": [
            "echo '{{message}}'"
        ]
    },
    "pipelines": {
        // Pipelines to be run
        // Keys are the pipeline names, values are the pipeline configurations
        "say-hello": {
            "servers": [
                "staging"
            ],
            "tasks": [
                "print-message"
            ]
        }
    }
}
```

**Auto populated params**:
- `{{timestamp}}` - The current timestamp.