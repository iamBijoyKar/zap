# ⚡ Zap

Zap is a powerful, lightweight task runner for projects that helps you automate and manage your development workflow. It executes tasks defined in a simple YAML configuration file with support for parallel execution, retries, dependencies, and detailed timing information.

## 🚀 Features

- **📋 Task Management**: Define and organize tasks in a simple YAML configuration
- **⚡ Parallel Execution**: Run multiple tasks concurrently for faster builds
- **🔄 Retry Mechanism**: Automatically retry failed tasks with configurable attempts
- **🔗 Dependency Management**: Define task dependencies to ensure proper execution order
- **⏱️ Execution Timing**: Track how long each task takes to complete
- **🔍 Verbose Output**: Get detailed information about task execution with `--verbose` flag
- **🛠️ Cross-Platform**: Works on Windows, macOS, and Linux
- **📦 Zero Dependencies**: Single binary with no external dependencies

## 📦 Installation

### Pre-built Binaries
Download the latest release from the [releases page](https://github.com/iamBijoyKar/zap/releases) and extract the binary to your PATH.

### Build from Source
```bash
git clone https://github.com/iamBijoyKar/zap.git
cd zap
go build -o zap
```

## 🎯 Quick Start

1. **Create a `zap.yaml` file** in your project root:
```yaml
tasks:
  - name: build
    command: ["go", "build", "-o", "app"]
    retries: 2

  - name: test
    command: ["go", "test", "./..."]
    depends_on: [build]

  - name: deploy
    command: ["docker", "build", "-t", "myapp", "."]
    parallel: true
    retries: 1
```

2. **Run all tasks**:
```bash
zap run all
```

3. **Run a specific task**:
```bash
zap run task build
```

4. **Enable verbose output**:
```bash
zap run all --verbose
```

## 📝 YAML Configuration Format

### Basic Structure
```yaml
tasks:
  - name: task-name
    command: ["command", "arg1", "arg2"]
    retries: 2
    parallel: true
    depends_on: [task1, task2]
```

### Task Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `name` | string | ✅ | Unique name for the task |
| `command` | array | ✅ | Command and arguments to execute |
| `retries` | integer | ❌ | Number of retry attempts (default: 0) |
| `parallel` | boolean | ❌ | Run task concurrently (default: false) |
| `depends_on` | array | ❌ | List of task names that must complete first |

### Example Configuration
```yaml
tasks:
  # Build task with retries
  - name: build
    command: ["npm", "run", "build"]
    retries: 3

  # Test task that depends on build
  - name: test
    command: ["npm", "test"]
    depends_on: [build]

  # Parallel deployment tasks
  - name: deploy-staging
    command: ["docker", "push", "myapp:staging"]
    parallel: true
    depends_on: [test]

  - name: deploy-production
    command: ["docker", "push", "myapp:latest"]
    parallel: true
    depends_on: [test]

  # Linting task (runs in parallel with other tasks)
  - name: lint
    command: ["eslint", "src/"]
    parallel: true
```

## 🎮 Usage

### Commands

#### Run All Tasks
```bash
zap run all [--verbose]
```
Executes all tasks in the configuration file, respecting dependencies and parallel execution.

#### Run Specific Task
```bash
zap run task <task-name> [--verbose]
```
Runs a specific task by name.

#### Help
```bash
zap --help
zap run --help
zap run all --help
zap run task --help
```

### Flags

| Flag | Description |
|------|-------------|
| `--verbose` | Enable detailed output including command execution details and retry information |

## 🔄 Execution Flow

1. **Parallel Tasks**: Tasks marked with `parallel: true` run concurrently
2. **Sequential Tasks**: Tasks without the parallel flag run in order
3. **Dependencies**: Tasks wait for their dependencies to complete successfully
4. **Retries**: Failed tasks are retried according to their `retries` configuration
5. **Timing**: Each task's execution time is tracked and displayed

## 📊 Example Output

### Basic Execution
```bash
$ zap run all

        ⚡ Zap 1.0.0 (golang)
         - Total Tasks: 4

1. Running Task ... 🔨 build
 > go build -o app
Task completed ✅ (took 1.2s)

2. Running Task ... 🔨 test
 > go test ./...
Task completed ✅ (took 3.4s)

3. Running Task ... 🔨 deploy
 > docker build -t myapp .
Task completed ✅ (took 15.6s)

Total Completed Tasks: 3
Total Failed Tasks: 0
```

### Verbose Execution
```bash
$ zap run all --verbose

        ⚡ Zap 1.0.0 (golang)
         - Total Tasks: 4

Running 2 parallel tasks concurrently...
Running parallel task: 🔨 deploy-staging
 > docker push myapp:staging
Running parallel task: 🔨 deploy-production
 > docker push myapp:latest
Executing command: docker push myapp:staging
Executing command: docker push myapp:latest
Parallel task 'deploy-staging' completed ✅ (took 8.2s)
Parallel task 'deploy-production' completed ✅ (took 9.1s)

1. Running Task ... 🔨 build
 > go build -o app
Executing command: go build -o app
Task completed ✅ (took 1.2s)

Total Completed Tasks: 3
Total Failed Tasks: 0
```

## 🛠️ Use Cases

- **CI/CD Pipelines**: Automate build, test, and deployment processes
- **Development Workflows**: Run linting, testing, and building tasks
- **Docker Workflows**: Build and push Docker images
- **Database Migrations**: Run database setup and migration tasks
- **Code Quality**: Execute linting, formatting, and security checks
- **Documentation**: Generate and deploy documentation

## 🔧 Advanced Configuration

### Complex Dependency Chain
```yaml
tasks:
  - name: clean
    command: ["rm", "-rf", "dist/"]
    
  - name: install-deps
    command: ["npm", "install"]
    depends_on: [clean]
    
  - name: build
    command: ["npm", "run", "build"]
    depends_on: [install-deps]
    retries: 2
    
  - name: test
    command: ["npm", "test"]
    depends_on: [build]
    
  - name: lint
    command: ["npm", "run", "lint"]
    parallel: true
    depends_on: [build]
    
  - name: deploy
    command: ["npm", "run", "deploy"]
    depends_on: [test, lint]
```

### Parallel Task Groups
```yaml
tasks:
  - name: build-frontend
    command: ["npm", "run", "build"]
    parallel: true
    
  - name: build-backend
    command: ["go", "build", "-o", "server"]
    parallel: true
    
  - name: build-docs
    command: ["mkdocs", "build"]
    parallel: true
    
  - name: package
    command: ["./package.sh"]
    depends_on: [build-frontend, build-backend, build-docs]
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

Zap is open-source software licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/)
- CLI framework powered by [urfave/cli](https://github.com/urfave/cli)
- YAML parsing with [go-yaml](https://github.com/go-yaml/yaml)
- Colored output with [fatih/color](https://github.com/fatih/color)