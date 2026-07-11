# Architecture

Entry point: `main` calls `shell.New(stdin, stdout, stderr).Run()`.

## Packages


| Package      | Responsibilities                                                                                      |
| ------------ | ----------------------------------------------------------------------------------------------------- |
| `shell`      | Runs the REPL loop, routes parsed input to execution, and coordinates jobs, history, and expansion.   |
| `terminal`   | Handles prompt display, raw-mode line editing, tab dispatch, history recall, and TTY output wrapping. |
| `parser`     | Parses input lines into commands, pipelines, redirects, and background flags.                         |
| `executor`   | Applies redirects and runs builtins, external programs, and pipelines.                                |
| `completer`  | Orchestrates tab completion for commands, arguments, filenames, and programmable scripts.             |
| `session`    | Holds mutable shell session state: jobs, history, variables, completion registry, and history file path. |
| `history`    | Maintains the in-memory command list and helpers to load, save, and format history entries            |
| `variables`  | Stores shell variables and expands parameter references in command arguments.                         |
| `jobs`       | Tracks background jobs and formats job listings for display.                                          |
| `completion` | Prefix-matches candidates, registers programmable completers, and runs completion scripts.            |
| `builtins`   | Implements and dispatches shell builtin commands.                                                     |
| `external`   | Resolves executables on PATH and runs external programs.                                              |
| `files`      | Lists directory entries for completion and provides line-oriented file I/O.                           |




## Dependency overview

```mermaid
flowchart TB
    main --> shell

    shell --> terminal
    shell --> executor
    shell --> parser
    shell --> completer
    shell --> session
    shell --> builtins
    shell --> external
    shell --> jobs
    shell --> variables

    session --> history
    session --> jobs
    session --> completion
    session --> variables

    history --> files

    completer --> builtins
    completer --> completion
    completer --> external
    completer --> files
    completer --> session
    completer --> terminal

    executor --> builtins
    executor --> parser
    executor --> external
    executor --> session

    builtins --> history
    builtins --> session
    builtins --> variables
```



## Class diagram

```mermaid
classDiagram
    direction TB

    class Shell {
        -terminal Terminal
        -executor Executor
        -completer Completer
        -state State
        +Run() error
        +ExecuteLine(string) bool, error
        +HandleTab(TabState, string) TabResult
        +HistoryPrevious(int) string, bool
    }

    class Terminal {
        -tabHandler TabHandler
        -historyHandler HistoryHandler
        -rawTTY RawMode
        -reader bufio.Reader
        +ReadLine() string, bool, error
        +WriteLine(string)
        +Stdout() io.Writer
        +Stderr() io.Writer
        +PrepareRead() bool
        +Close() error
    }

    class HistoryHandler {
        <<interface>>
        +HistoryPrevious(int) string, bool
    }

    class RawMode {
        +PrepareRead() bool
        +Active() bool
        +Close() error
    }

    class TabHandler {
        <<interface>>
        +HandleTab(TabState, string) TabResult
    }

    class TabState {
        +PendingListings []string
    }

    class TabResult {
        +Buffer string
        +ListingsToShow []string
        +RingBell bool
    }

    class Completer {
        -state State
        +HandleTab(TabState, string) TabResult
        +ApplyTabAction(TabState, string, string, []string) TabResult
    }

    class Executor {
        -stdin io.Reader
        +ExecuteBuiltin(Outputs, State, []string) bool, error
        +ExecuteExternalForeground(Outputs, []string) error
        +ExecuteExternalBackground(Outputs, []string, func()) int, error
        +ExecutePipeline(Outputs, State, [][]string) error
    }

    class Outputs {
        +Stdout io.Writer
        +Stderr io.Writer
        +Redirect Redirect
    }

    class State {
        +Jobs Table
        +History List
        +Histfile string
        +Completion Registry
        +Variables Store
    }

    class List {
        +Add(string)
        +List() []Entry
        +ListLast(int) []Entry
        +Previous(int) string, bool
        +ReadFromFile(string) error
        +AppendFromFile(string) error
        +WriteToFile(string) error
        +AppendToFile(string) error
    }

    class Entry {
        +Number int
        +Command string
    }

    class Table {
        +Add(int, string) int
        +MarkDone(int)
        +ReapDone() []Job
        +List() []Job
    }

    class Job {
        +Number int
        +PID int
        +Command string
        +Status string
    }

    class Registry {
        +Register(string, string)
        +Unregister(string)
        +Lookup(string) string, bool
    }

    class Store {
        +Set(string, string)
        +Get(string) string, bool
    }

    class variables {
        <<package>>
        +IsValidIdentifier(string) bool
        +ExpandField(Store, string) string
        +ExpandFields(Store, []string) []string
    }

    class builtinsRegistry {
        +Register(string, Handler)
        +Run(string, []string, Context) bool, error
        +Is(string) bool
        +Names() []string
    }

    class Context {
        +Stdout io.Writer
        +Stderr io.Writer
        +State State
    }

    class builtins {
        <<package>>
        +IsBuiltin(string) bool
        +Names() []string
        +Run(string, []string, Context) bool, error
    }

    class parser {
        <<package>>
        +Tokenize(string) []string
        +ParseLine(string) Line
        +ParseCommand([]string) Command
        +ParsePipelineSegments([][]string) [][]string, Redirect
        +SplitPipelineTokens([]string) [][]string
        +ParseRedirect([]string) []string, Redirect
        +StripBackground([]string) []string, bool
    }

    class Line {
        +Pipeline bool
        +Commands [][]string
        +Redirect Redirect
        +Background bool
    }

    class Redirect {
        +StdoutPath string
        +StdoutAppend bool
        +StderrPath string
        +StderrAppend bool
    }

    class external {
        <<package>>
        +FindExecutableInPath(string) string, bool
        +FindAllExecutablesInPath() []string
        +New([]string, io.Writer, io.Writer) ExternalProgram, bool
    }

    class ExternalProgram {
        +Run() error
        +RunInBackground(func()) int, error
    }

    class completion {
        <<package>>
        +Complete(string, []string) string, []string, bool
        +RunCompleter(CompleterOptions) []string, error
    }

    class jobs {
        <<package>>
        +FormatLines([]Job) []string
        +WriteAll(io.Writer, []Job)
    }

    class files {
        <<package>>
        +ListInDir(string, string) []string
        +ReadLines(string) []string, error
        +WriteLines(string, []string) error
        +AppendLines(string, []string) error
    }

    class history {
        <<package>>
        +FormatLines([]Entry) []string
        +WriteAll(io.Writer, []Entry)
    }

    Shell --> Terminal
    Shell --> Executor
    Shell --> Completer
    Shell --> State
    Shell ..|> TabHandler
    Shell ..|> HistoryHandler
    Shell ..> parser
    Shell ..> builtins
    Shell ..> external
    Shell ..> jobs

    Terminal --> TabHandler
    Terminal --> HistoryHandler
    Terminal --> RawMode
    TabHandler ..> TabState
    TabHandler ..> TabResult

    Completer --> State
    State --> Table
    State --> List
    State --> Registry
    State --> Store
    Table --> Job
    List --> Entry

    Executor --> Outputs
    Outputs --> Redirect
    Executor ..> builtins
    Executor ..> external
    Executor ..> State

    builtins ..> builtinsRegistry
    builtins ..> Context
    Context --> State
    builtins ..> variables
    Shell ..> variables

    parser ..> Line
    parser ..> Redirect
    external ..> ExternalProgram

    Completer ..> completion
    Completer ..> builtins
    Completer ..> external
    Completer ..> files
    history ..> files
```





## REPL loop

Owned by `Shell.Run()`:

1. `terminal.PrepareRead()` — re-enable raw mode (external programs may restore cooked mode)
2. `writeReapedJobs()` — `state.Jobs.ReapDone()` → `jobs.FormatLines` → `terminal.WriteLine` each line
3. `terminal.ReadLine()` — up/down arrows browse history via `HistoryHandler` when raw mode is active
4. `ExecuteLine(line)` — `state.History.Add(line)`, then `parser.ParseLine`, `expandParsedLine` (`variables.ExpandFields`), `commandFound`, dispatch to `executor`
5. Repeat until exit or EOF

On exit, a deferred handler in `Run()` writes `state.History` to `state.Histfile` when `HISTFILE` was set at startup.

`ExecuteLine` branches on `parsed.Pipeline`: single commands go to `executeCommand`; pipelines go to `executePipeline` (which validates every segment via `validatePipelineSegments`).

## Tab completion

Owned by `completer` package; `Shell.HandleTab` delegates to `completer.Completer`.


| File           | Role                                                                                         |
| -------------- | -------------------------------------------------------------------------------------------- |
| `completer.go` | Routing, programmable completion (`BuildCompleterOptions`), double-Tab UX (`ApplyTabAction`) |
| `command.go`   | First-token completion: `builtins.Names` + PATH executables                                  |
| `file.go`      | Filename candidate sourcing via `files.ListInDir`                                            |
| `argument.go`  | Last-argument prefix matching (files and programmable candidates)                            |


Flow:

1. User presses Tab during `terminal.ReadLine()`
2. `terminal` calls `tabHandler.HandleTab(state, buffer)` — implemented by `Shell`
3. `Completer.completeBuffer` routes to command, programmable, or filename completion
4. `completion.Complete` runs prefix matching on gathered candidates
5. `ApplyTabAction` applies double-Tab logic (bell on first Tab, listings on second)
6. `terminal` updates the buffer or shows match listings

The `complete` builtin registers and unregisters scripts via `session.State.Completion` (`completion.Registry`).

## Terminal I/O

- **Raw mode** (`terminal/raw.go`, type `RawMode`): byte-at-a-time input so Tab, Backspace, arrow-key history, and completion listings work. Falls back to line-based reads when stdin is not a TTY (tests).
- **Command writers** (`terminal.Stdout()` / `Stderr()`): called at execution time, not cached. When raw mode is active, `WrapWriter` (`writer.go`) translates `\n` → `\r\n` so each line starts at column 0.
- **Input**: `bufio.Reader` on stdin for `ReadLine`; `RawMode` holds the `*os.File` for `MakeRaw` / `Restore`.
- **Tab types** (`terminal/tab.go`): `TabHandler` interface, `TabState`, `TabResult` — keeps completion semantics out of `terminal`.
- **History recall** (`terminal/history.go`, `input.go`): `HistoryHandler` interface; `historyBrowseState` tracks up/down navigation on the current prompt line. `handleEscapeSequence` handles CSI arrow sequences (`ESC [ A` / `ESC [ B`); `Shell` implements `HistoryPrevious` by delegating to `state.History.Previous`. Manual edits reset browse position.



## Parsing

`parser.ParseLine` (`parser/parser.go`) is the single entry point for line parsing:

- Single command → `ParseCommand` (`ParseRedirect` + `StripBackground`)
- Pipeline → `SplitPipelineTokens` + `ParsePipelineSegments` (redirect on final segment; background `&` stripped)



## Executor

Public API lives in `executor.go`. Private stage runners live in `run.go`; pipeline wiring in `pipeline.go`; redirect open/close in `redirect.go`.


| Method                      | Role                                                                          |
| --------------------------- | ----------------------------------------------------------------------------- |
| `ExecuteBuiltin`            | `withOutputs` → `runBuiltin` → `builtins.Run`                                 |
| `ExecuteExternalForeground` | `withOutputs` → `runExternal` (stdin from executor)                           |
| `ExecuteExternalBackground` | `withOutputs` → `runExternalBackground`; returns PID only                     |
| `ExecutePipeline`           | `withOutputs` → `runPipeline` (goroutine per stage, `io.Pipe` between stages) |


Shared private runners in `run.go`:

- `runBuiltin` — builds `builtins.Context`; drains pipe stdin for middle pipeline builtins via `runDrainingStdin`
- `runExternal` — `external.New` + `Run`
- `runExternalBackground` — `RunInBackground` with caller-supplied `onExit` callback
- `nonExitError` — swallows `exec.ExitError` for foreground commands

`Shell` builds `executor.Outputs` on each command/pipeline from `terminal.Stdout()`, `terminal.Stderr()`, and the parsed redirect. `session.State` is passed per call for builtins and pipeline stages that need jobs/completion.

## Builtin commands

`builtins` package holds implementations and a unified registry (`registry.go`) patterned after `completion/registry.go`. Each command self-registers via unexported `register()` in its file's `init()`.

Registered builtins: `cd`, `complete`, `declare`, `echo`, `exit`, `history`, `jobs`, `pwd`, `type`.


| Concern                                                    | Owner                                                           |
| ---------------------------------------------------------- | --------------------------------------------------------------- |
| Builtin implementations                                    | Per-command files (`echo.go`, `cd.go`, `exit.go`, …)            |
| Registry (`builtins.Registry`)                             | `builtins/registry.go` — `Run`, `IsBuiltin`, `Names`            |
| Shell session state (jobs, history, completion, variables) | `session.State`, owned by `Shell`                                  |
| Per-invocation I/O and state refs                          | `builtins.Context` (`Stdout`, `Stderr`, `State *session.State`)    |
| Invoking builtins                                          | `Executor` → `builtins.Run`                                     |
| Routing builtin vs external                                | `Shell.ExecuteLine` via `commandFound` and `builtins.IsBuiltin` |
| Command resolution (`type`, pre-exec check)                | `builtins/type.go` (`Type`) and `shell.commandFound`            |




## Command resolution and shell messages

`commandFound` (package-private helper in `shell.go`) checks `builtins.IsBuiltin` then `external.FindExecutableInPath` before execution. The `type` builtin uses the same classification inline in `Type` with different message formatting.

`Shell.ExecuteLine` resolves the command before calling executor:

- `commandFound` → `CommandNotFoundMessage` via `terminal.WriteLine`, or continue
- `builtins.IsBuiltin` → `ExecuteBuiltin`
- `parsed.Background` → `executeBackgroundCommand`
- otherwise → `ExecuteExternalForeground`

Background jobs: `executeBackgroundCommand` starts the process via `ExecuteExternalBackground`, registers the job in `state.Jobs` (`Add`, `MarkDone` callback), and prints `[n] pid`. Reaped jobs print before the next prompt via `writeReapedJobs`.

## Parameter expansion and shell variables

Shell variables live in `variables.Store` on `session.State`. The `declare` builtin creates and inspects them; `shell.ExecuteLine` expands `$VAR` and `${VAR}` references in parsed command arguments before execution.


| Concern                      | Owner                                                             |
| ---------------------------- | ----------------------------------------------------------------- |
| In-memory variable store     | `variables.Store` on `session.State`                                 |
| Create / inspect variables   | `builtins/declare.go` — assignment form and `declare -p`          |
| Identifier validation        | `variables.IsValidIdentifier` (shared by `declare` and expansion) |
| Expand parsed command fields | `shell.expandParsedLine` → `variables.ExpandFields`               |
| Per-field `$VAR` / `${VAR}`  | `variables.ExpandField` (`variables/expand.go`)                   |




### Expansion timing

Expansion runs after `parser.ParseLine` and before `commandFound`. Each token in every pipeline segment is expanded independently. The original command line stored in history is unexpanded.

`variables.ExpandFields` drops empty arguments after expansion. For example, `${missing}` alone becomes an empty word and is omitted from the argv passed to external programs.

### Supported forms


| Form     | Example          | Result when `foo=bar` | Result when unset                   |
| -------- | ---------------- | --------------------- | ----------------------------------- |
| `$VAR`   | `echo $foo`      | `echo bar`            | `echo` (empty arg dropped if alone) |
| `${VAR}` | `echo ${foo}end` | `echo barend`         | `echo end`                          |
| Adjacent | `echo $foo$foo`  | `echo barbar`         | `echo`                              |


Invalid identifiers (for example `$1foo` or `${1foo}`) are left literal. A bare `$` with no valid name is also left literal.

### `declare` builtin


| Invocation           | Behavior                                                                    |
| -------------------- | --------------------------------------------------------------------------- |
| `declare name=value` | Parse `name=value`, validate identifier, `store.Set(name, value)`           |
| `declare -p name`    | Print `declare -- name="value"` or `declare: name: not found` on stderr     |
| Invalid name         | `declare: \`name=value': not a valid identifier` on stderr; no store update |


Assignment parsing (`parseAssignment`) lives in `builtins/declare.go`. Store operations and expansion logic stay in `variables`.

### Expansion flow

```mermaid
sequenceDiagram
    participant shell as Shell.ExecuteLine
    participant parser as parser.ParseLine
    participant vars as variables.ExpandFields
    participant exec as executor

    shell->>parser: ParseLine(line)
    parser-->>shell: Line with tokenized Commands
    shell->>vars: ExpandFields(state.Variables, fields)
    vars-->>shell: expanded argv (empty words dropped)
    shell->>shell: commandFound(expanded fields)
    shell->>exec: ExecuteBuiltin / ExecuteExternal / ExecutePipeline
```





## Command history

The `history` package is path-agnostic: it stores commands in memory and exposes file helpers that take an explicit path. `HISTFILE` policy (read on startup, write on exit) lives in `session` and `shell`; the builtin supplies paths for `-r`, `-w`, and `-a`.


| Concern                   | Owner                                                                  |
| ------------------------- | ---------------------------------------------------------------------- |
| In-memory command list    | `history.List` on `session.State`                                         |
| `HISTFILE` path           | `session.State.Histfile` from `os.Getenv("HISTFILE")` in `session.NewState`  |
| Load history on startup   | `session.NewState` → `History.AppendFromFile(histfile)` (missing file OK) |
| Record each executed line | `shell.ExecuteLine` → `state.History.Add(line)` before parsing         |
| Persist on shell exit     | `shell.Run` defer → `state.History.WriteToFile(state.Histfile)`        |
| List / file ops builtin   | `builtins/history.go` → `history.List` methods                         |
| Up/down arrow recall      | `terminal/history.go` + `input.go`; `Shell.HistoryPrevious`            |
| Line-oriented file I/O    | `files.ReadLines`, `files.WriteLines`, `files.AppendLines`             |




### `history` package

`List` (`history/history.go`) is a mutex-protected slice of command strings.


| Method              | Role                                                                                        |
| ------------------- | ------------------------------------------------------------------------------------------- |
| `Add`               | Append a command after the user submits a line                                              |
| `List` / `ListLast` | Snapshot entries with bash-style line numbers (`Entry.Number` preserves original index)     |
| `Previous`          | Random access for arrow recall (`stepsBack` 0 = most recent)                                |
| `ReadFromFile`      | Append lines from a path; errors propagate (used by `history -r`)                           |
| `AppendFromFile`    | Append lines; empty path and missing file are no-ops (used for `HISTFILE` load)             |
| `WriteToFile`       | Overwrite a path with the full list (used for `HISTFILE` exit and `history -w`)             |
| `AppendToFile`      | Append commands since the last file read/write/append (`history -a`; tracks `lastAppended`) |


Display helpers `FormatLines` and `WriteAll` produce bash-style output (`%5d  %s` per entry). The builtin delegates listing to `builtins.History`, which chooses `List` vs `ListLast` from an optional numeric limit.

### `history` builtin


| Invocation        | Behavior                                                             |
| ----------------- | -------------------------------------------------------------------- |
| `history`         | Print full history via `history.WriteAll`                            |
| `history n`       | Print last *n* entries (`n` must be a positive integer)              |
| `history -r path` | `ReadFromFile(path)` — append file contents; errors to stderr        |
| `history -w path` | `WriteToFile(path)` — overwrite file with full list                  |
| `history -a path` | `AppendToFile(path)` — append only commands added since last file op |




### Persistence flow

```mermaid
sequenceDiagram
    participant session as session.NewState
    participant hist as history.List
    participant shell as Shell.Run
    participant term as terminal.ReadLine

    session->>hist: AppendFromFile(HISTFILE)
    loop REPL
        term->>shell: ReadLine (arrow keys via HistoryPrevious)
        shell->>hist: Add(line)
        shell->>shell: ExecuteLine
    end
    shell->>hist: WriteToFile(HISTFILE) on exit
```





## Background jobs


| Concern              | Owner                                        |
| -------------------- | -------------------------------------------- |
| Job table data       | `jobs.Table` on `session.State`                 |
| Start + `[n] pid` UX | `shell.executeBackgroundCommand`             |
| Mark done on exit    | `onExit` callback from executor → `shell`    |
| Reap + print         | `shell.writeReapedJobs` → `jobs.FormatLines` |
| List on demand       | `builtins/jobs` → `jobs.WriteAll`            |


