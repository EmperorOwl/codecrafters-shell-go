[![progress-banner](https://backend.codecrafters.io/progress/shell/ddc43dfa-46b2-40f0-9cf8-439ddefa7004)](https://app.codecrafters.io/users/EmperorOwl?r=2qF)

# Shell

A POSIX-style interactive shell built in Go for the [CodeCrafters Shell challenge](https://app.codecrafters.io/courses/shell/overview). The implementation parses user input, runs builtins and external programs, supports pipelines and redirects, and provides tab completion with optional programmable scripts.

## Features

- Interactive REPL with raw-mode line editing, history recall, and tab completion
- Builtin commands: `echo`, `exit`, `pwd`, `cd`, `type`, `declare`, `complete`, `jobs`, `history`
- External command execution with PATH lookup
- Pipelines, stdout/stderr redirects (including append), and background jobs
- Shell variables with `$VAR` / `${VAR}` expansion
- Persistent command history via `HISTFILE`
- Programmable tab completion via `complete -C <script> <command>`

## Requirements

- Go 1.26+

## Running locally

```sh
./your_program.sh
```

The entry point is `app/main.go`, which constructs a `shell.Shell` and calls `Run()`.

## Testing

Run the full test suite:

```sh
go test ./...
```

Run tests with coverage, excluding the `testutils` package (shared test helpers, not production code):

```sh
go test $(go list ./... | grep -v testutils) -cover
```

Generate a coverage report:

```sh
go test $(go list ./... | grep -v testutils) -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### Coverage

Overall statement coverage is **~85%** across production packages (`testutils` excluded).

| Package     | Coverage |
| ----------- | -------- |
| `parser`    | 100%     |
| `session`   | 100%     |
| `completer` | 92%      |
| `files`     | 91%      |
| `history`   | 91%      |
| `variables` | 93%      |
| `completion`| 90%      |
| `executor`  | 90%      |
| `builtins`  | 87%      |
| `jobs`      | 89%      |
| `shell`     | 86%      |
| `external`  | 82%      |
| `terminal`  | 57%      |

`terminal` has the lowest coverage because interactive TTY input paths are harder to exercise in unit tests.

## Architecture

See [ARCHITECTURE.MD](ARCHITECTURE.MD) for package responsibilities, dependency diagrams, and class relationships.

## CodeCrafters

Submit your solution to the CodeCrafters servers:

```sh
codecrafters submit
```

If you're viewing this repo on GitHub, head over to [codecrafters.io](https://codecrafters.io) to try the challenge.
