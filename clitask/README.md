# CLI Task Manager

## Exercise details

In this exercise we are going to be building a CLI tool that can be used to manage your TODOs in the terminal. The basic usage of the tool is going to look roughly like this:

```
$ task
task is a CLI for managing your TODOs.

Usage:
  task [command]

Available Commands:
  add         Add a new task to your TODO list
  do          Mark a task on your TODO list as complete
  list        List all of your incomplete tasks

Use "task [command] --help" for more information about a command.

$ task add review talk proposal
Added "review talk proposal" to your task list.

$ task add clean dishes
Added "clean dishes" to your task list.

$ task list
You have the following tasks:
1. review talk proposal
2. some task description

$ task do 1
You have completed the "review talk proposal" task.

$ task list
You have the following tasks:
1. some task description
```

*Note: Lines prefixed with `$` are lines where we type into the terminal, and other lines are output from our program.*

### 1. Build the CLI shell

Your final CLI won't need to look exactly like this, but this is what I roughly expect mine to look like. In the bonus section we will also discuss a few extra features we could add, but for now we will stick with the three show above:

- `add` - adds a new task to our list
- `list` - lists all of our incomplete tasks
- `do` - marks a task as complete

### 2. Write the BoltDB interactions

After stubbing out your CLI commands, try writing code that will read, add, and delete data in a BoltDB database. You can find more information about using Bolt here: <https://github.com/etcd-io/bbolt>

For now, don't worry about where you store the database that bolt connects to. At this stage I intend to just use whatever directory the `task` command was run from, so I will be using code roughly like this:

```go
db, err := bolt.Open("tasks.db", 0600, nil)
```

Later you can dig into how to install the application so that it can be run from anywhere and it will persist our tasks regardless of where we run the CLI.
