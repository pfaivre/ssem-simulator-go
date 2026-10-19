# AGENTS.md

Go simulator of the SSEM (Manchester Baby, 1948), the first stored-program computer. Standard library only, no external dependencies. Module: `pfaivre/ssem-simulator-go`, Go 1.24.

## Setup

Any command should be run inside the devenv like so:

```sh
nix shell nixpkgs#devenv -c devenv shell <command to run>
```

## Commands

```sh
go build ./...                          # Build everything
go test ./...                           # Run all tests (also run by devenv on enter)
go test ./... -cover                    # Coverage
go test ./... -race                     # Race condition detection
go vet ./...                            # Static checks
gofmt -l .                              # Formatting check (should output nothing)
go run ./main.go --file ./samples/ssem/nightmare.snp          # Run with a binary snapshot
go run ./main.go --file ./samples/ssem/fibonacci.asm          # Run with assembly
```

CLI flags (main.go): `--file` (.asm or .snp), `--pretty` (default true, maps 0/1 to ./#), `--speed` (cycles per second, 0 = unlimited), `--max-cycles`.

Development environment is a Nix devenv (`devenv.nix`); only Go is required to build and test.

## Architecture

- `main.go` — CLI entry point. Parses flags, loads a program, then runs two goroutines: `machine.Run` (compute) and `machine.Printer` (display, ~60 fps), coordinated with channels. Handles SIGINT/SIGTERM.
- `ssem/` — the machine itself, no CLI concerns:
  - `simulator.go` — `Ssem` state (store, CI, accumulator, `StopFlag`), file loaders (`ReadAsm`, `ReadSnp`), fetch-decode-execute cycle (`InstructionCycle`, `Execute`), `Run` and `Printer`.
  - `store.go` — `Word` (int32), `Store` ([32]Word), instruction decoding (opcode bits 13-15, address bits 0-4), binary rendering with ANSI colors.
  - `opcode.go` — `Opcode` type, mnemonics (JMP, JRP, LDN, STO, SUB, SUB2, CMP, STP, plus pseudo-mnemonic NUM for raw data), parsing tables.

## Domain rules that are easy to get wrong

- The SSEM stores integers **least significant bit first**. On-disk .snp lines and `fmt.Sprintf("%032b")` output are MSB-first, so `Reverse` is applied at every boundary (load, print).
- `NUM` (0b1111) is not a real instruction; it is an assembly-only pseudo-mnemonic that stores a raw value. When the word is decoded as an instruction, bits 13-15 are masked to 3 bits regardless.
- Addresses are 5 bits, words 32 bits, store is 32 words. CI wraps modulo 32.
- CMP historically uses a Test Flip-Flop; the current implementation approximates it by skipping only the next instruction. Accuracy validation against the reference literature is an open roadmap item (see README).

## Conventions

- Keep changes inside `ssem/` free of `os`, `flag`, or anything CLI-specific; the package should stay importable as a plain simulator.
- Error handling: loaders return descriptive errors with line numbers (`line %d: ...`), printed to stderr by main.
- Tests live next to the code (`ssem/simulator_test.go`, `ssem/store_test.go`), including a functional test that runs a full program from a store snapshot and compares the final store. Use `storeFromSnapshot` for store fixtures; snapshot strings are MSB-first like .snp files.
- Sample programs: `samples/ssem/*.asm` (assembly) and `*.snp` (binary), `samples/ssem/tests/` holds per-instruction test programs.
- Comments and identifiers are in English; watch the existing style: constants are SCREAMING_SNAKE, exported functions get doc comments.
