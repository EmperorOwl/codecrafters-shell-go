# Architecture

Entry point: `main` calls `shell.New(stdin, stdout, stderr).Run()`.

| Package | Responsibilities |
| --- | --- |
| `shell` | Top-level orchestrator: owns the REPL loop (`shell.go`), tab completion (`tab.go`), terminal, executor, job state, and completion registry. |
| `terminal` | Handles user I/O: prompt, line editing, raw mode, Tab key handling, and command output writers. |
| `parser` | Tokenizes and parses input into commands, arguments, pipelines, and redirects. |
| `executor` | Opens redirect outputs and runs commands: builtins, external programs, and pipelines. |
| `jobs` | Tracks background jobs: add, mark done, reap, and list. |
| `completion` | Tab match logic (`Complete`), programmable completion registry (`CompletionRegistry`), and completer script execution (`RunCompleter`). |
| `builtins` | Builtin command implementations and dispatch (`echo`, `cd`, `exit`, `type`, `jobs`, `complete`, `pwd`). |
| `external` | PATH lookup and running external programs (`ExternalProgram`). |
| `files` | Directory listing for file tab completion. |

## API style

| Style | Packages / types | Reason |
| --- | --- | --- |
| **Struct** (owned state or injected deps) | `Shell`, `Terminal`, `Executor`, `JobTable`, `CompletionRegistry`, `ExternalProgram` | Session state, lifecycle, or dependencies wired at `New()` |
| **Package functions** (stateless) | `parser`, `completion`, `builtins`, `files`; `external` PATH helpers | Pure input→output; no per-shell instance needed |
| **Types only** | `Redirect`, `Job`, `State`, `BuiltinContext`, `CompleterOptions`, `TabState`, `TabResult` | Data passed between layers; `CompleterOptions` lives in `completion` |

## Class diagram

```mermaid
classDiagram
    class Shell {
        -terminal Terminal
        -executor Executor
        -state State
        +Run() error
        +ExecuteLine(string) bool, error
        +HandleTab(TabState, string) TabResult
    }

    class parser {
        <<package>>
        +Tokenize(string) []string
        +SplitPipelineTokens([]string) [][]string
        +StripBackground([]string) []string, bool
        +ParseRedirect([]string) []string, Redirect
    }

    class Redirect {
        +StdoutPath string
        +StdoutAppend bool
        +StderrPath string
        +StderrAppend bool
    }

    parser ..> Redirect

    class Outputs {
        +Stdout io.Writer
        +Stderr io.Writer
        +Redirect Redirect
    }

    class Executor {
        -stdin io.Reader
        +New(io.Reader)
        +ExecuteBuiltin(Outputs, State, []string) bool, error
        +ExecuteExternalForeground(Outputs, []string) error
        +ExecuteExternalBackground(Outputs, []string, func()) int, error
        +ExecutePipeline(Outputs, State, [][]string) error
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

    external ..> ExternalProgram

    class builtins {
        <<package>>
        +IsBuiltin(string) bool
        +Names() []string
        +Run(string, []string, BuiltinContext) bool, error
    }

    class State {
        +Jobs JobTable
        +Completion CompletionRegistry
    }

    class BuiltinContext {
        +Stdout io.Writer
        +Stderr io.Writer
        +State State
    }

    builtins ..> State
    builtins ..> BuiltinContext

    class Job {
        +Number int
        +PID int
        +Command string
        +Status string
    }

    class JobTable {
        +Add(int, string) int
        +MarkDone(int)
        +ReapDone() []Job
        +List() []Job
    }

    class jobs {
        <<package>>
        +FormatLines([]Job) []string
    }

    JobTable ..> Job
    jobs ..> Job

    class completion {
        <<package>>
        +Complete(string, []string) string, []string, bool
        +RunCompleter(CompleterOptions) []string, error
    }

    class CompletionRegistry {
        +Register(string, string)
        +Unregister(string)
        +Lookup(string) string, bool
    }

    class CompleterOptions {
        +Path string
        +Command string
        +CurrentWord string
        +PreviousWord string
        +CompLine string
        +CompPoint int
    }

    completion ..> CompletionRegistry
    completion ..> CompleterOptions

    class files {
        <<package>>
        +ListInDir(string, string) []string
    }

    class TabHandler {
        <<interface>>
        +HandleTab(TabState, string) TabResult
    }

    class TabState {
        +pendingListings []string
    }

    class TabResult {
        +Buffer string
        +ListingsToShow []string
        +RingBell bool
    }

    class Terminal {
        -tabHandler TabHandler
        +ReadLine() string, bool, error
        +WriteLine(string)
        +Stdout() io.Writer
        +Stderr() io.Writer
    }

    TabHandler ..> TabState
    TabHandler ..> TabResult

    Executor ..> Redirect
    Executor ..> external
    Executor ..> builtins
    Shell --> Terminal : terminal
    Shell --> Executor : executor
    Shell --> State : state
    Shell ..|> TabHandler : implements
    Shell ..> parser
    Shell ..> external
    Shell ..> completion
    Shell ..> files
    Shell ..> builtins
    Shell ..> jobs
    Terminal --> TabHandler : tabHandler
    Executor ..> BuiltinContext : builds per call
```

## REPL loop

Owned by `Shell.Run()`:

1. `writeReapedJobs()` — `state.Jobs.ReapDone()` → `jobs.FormatLines` → `terminal.WriteLine` each line
2. `terminal.ReadLine()`
3. `ExecuteLine(line)` — `parser.*`, resolve command, dispatch to `executor` (redirect open/close handled inside executor)
4. Repeat until exit or EOF

## Tab completion

Owned by `shell/tab.go`.

1. User presses Tab during `terminal.ReadLine()`
2. `terminal` calls `tabHandler.HandleTab(state, buffer)` — implemented by `Shell`
3. `Shell.completeBuffer` routes to command, programmable, or filename completion:
   - **Commands:** deduplicated `builtins.Names` + PATH (`commandCandidates`)
   - **Programmable:** `buildCompleterOptions` → `state.Completion.Lookup` → `completion.RunCompleter`
   - **Files:** `files.ListInDir` for the current argument
4. `Shell` calls `completion.Complete` on the gathered candidates
5. `Shell` applies double-Tab logic (bell on first Tab, listings on second) and returns `TabResult`
6. `terminal` updates the buffer or shows match listings

The `complete` builtin registers and unregisters scripts via `CompletionRegistry`.

## Builtin commands

`builtins` package holds implementations and a fixed handler table (package-level `IsBuiltin`, `Names`, `Run`).

| Concern | Owner |
| --- | --- |
| Builtin implementations (`echo`, `cd`, …) | `builtins` package |
| Dispatch (`Run`, `IsBuiltin`, `Names`) | `builtins` package functions |
| Per-invocation I/O and shell state | `BuiltinContext` with `State` |
| Invoking builtins (single command or pipeline stage) | `Executor.ExecuteBuiltin` → `builtins.Run` |
| Routing builtin vs external | `Shell.ExecuteLine` via `builtins.IsBuiltin` and `external.FindExecutableInPath` |
| Builtin names for tab completion | `shell/tab.go` via `commandCandidates` (`builtins.Names` + PATH) |

`Shell` owns `builtins.State` (jobs and completion registry) and passes `executor.Outputs` plus `State` into execute methods at execution time so raw-mode LF translation is resolved when commands run. `Executor` stores stdin from `New`, starts processes, and builds `BuiltinContext` per builtin call. Background job registration (`Add`, `MarkDone`, `[n] pid` output) is handled by `Shell`.

Individual builtins stay as testable functions (e.g. `Echo`, `Cd`, `Type`) with thin handlers registered in the handler table.

## Command resolution and shell messages

`Shell.ExecuteLine` resolves the command before calling executor:

- `builtins.IsBuiltin` → `ExecuteBuiltin`
- `external.FindExecutableInPath` → foreground or background external execution
- Neither → `terminal.WriteLine` with command-not-found message

Background job startup (`[n] pid`) is printed by `Shell` after `ExecuteExternalBackground` returns a PID and `JobTable.Add` assigns a job number.
