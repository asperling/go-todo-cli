# 🐹 todo

A simple command-line todo manager written in Go for those who have at least one open terminal at any point in time.

This project started as a way to learn Go — its syntax, tooling, and ecosystem. I now actually use it daily, since it's lightweight, lives in the terminal, and is always just a few keystrokes away. Not as a todo manager for long-lived tasks, but rather to streamline my thoughts when working on, for example, a ticket. Usually, my commits reflect the todos I had set up beforehand.

## ✨ Features

- Add, list, complete, move, and delete todos
- Multiple task lists ("packages") via `todo packages`
- Configurable storage path
- Minimal dependencies, fast and local

```
 💋  ~/Code $ todo
NAME:
   todo - Manage your todos from the command line

USAGE:
   todo [global options] command [command options]

VERSION:
   0.1.0

COMMANDS:
   init, i           Initialize the todo CLI configuration
   list, ls, l       List all todos in the active package
   add, a            Add a new task to your current package
   delete, del, rm   Delete a task by its index
   done, d           Mark a task as completed
   undone, u         Reopen a completed task
   move, mv, m       Change the position of a task in your current package
   packages, pkg, p  Manage todo packages. See `todo packages help` for more information.
   help, h           Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## 🛠 Installation

```bash
git clone https://github.com/asperling/go-todo-cli.git
cd go-todo-cli
make install
```

Requires Go 1.24 or newer.

## 🚀 Getting Started

You could point your storage path to one that gets synched by e.g. Google Drive
to be able to access your todos on any machine.

```bash
todo init       # Setup config and storage path
todo add "Learn Go"
todo list
```

## 🔍 Why Go?

This project helped me explore:

- the Go module system and standard library
- structuring CLI apps with urfave/cli/v2
- unit testing and temporary file handling
- writing tools I actually use

## License

MIT
