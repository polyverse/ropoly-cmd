# polyverse/ropoly-cmd

## Build instructions for Docker

`build.sh` produces both a binary and a Docker image each called `ropoly-cmd`. The Docker image has the `ropoly-cmd`
binary and a directory `TestFiles` containing some simple fake fingerprints for testing the EQI calculation.

## Usage

`ropoly-cmd <input format> <output format> [additional args needed depending on output format]`

The behavior depends on the output format, but generally `ropoly-cmd` takes input from a file in the format specified
by `<input format>`, converts it to the format specified by `<output format>`, and prints the result to the command line.

For example, the output format `eqi` requires two inputs representing files and produces an EQI score between 0 and 100
inclusive rating how much the second binary's gadgets differ from the first. To get a score rating a scrambled ELF file's
gadgets' difference from the original file it is based on, use `ropoly-cmd file eqi path/to/original path/to/scrambled`.

### Input/output formats

#### file

Only valid as an input format, not output.

As an input, represents a filepath to an ELF binary file. (PE binaries are not yet supported.)

#### pid

Not yet implemented.

Only valid as an input format, not output.

As an input, represents a PID (not a filepath; exception to the usual rule) of a currently running ELF program.
(PE binaries are not yet supported.)

#### bindump

A JSON object containing the executable segments of a binary.

As an input, represents a filepath to such a JSON object.

As the output type, the JSON object is printed to stdout.

##### Usage for `bindump` output

`ropoly-cmd <input format> bindump <path/to/input/file>`

Valid input formats are `file`, `pid`, and `bindump` (conversion from `bindump` to `bindump` is a no-op).

#### fingerprint

A JSON object containing gadgets extracted from a binary.

As an input, represents a filepath to such a JSON object.

As an output type, the JSON object is printed to stdout.

##### Usage for `fingerpint` output

`ropoly-cmd <input format> fingerprint <path/to/input/file> [min-gadget-length <non-negative integer>] [max-gadget-length <non-negative integer>]`

Valid input formats are `file`, `pid`, `bindump`, and `fingerprint` (conversion from `fingerpint` to `fingerprint` is a no-op).

Optionally, you can specify the lengths of gadgets to be included in the fingerprint. See "Specifying gadget length" for details.

#### eqi

Only valid as an output format, not input.

A number between 0 and 100 inclusive representing the difference in gadgets between a modified binary and an original binary,
with 0 being the least different and 100 meaning that the two binaries share no gadgets whatsoever.

As the output type, the EQI score is printed to stdout.

##### Usage for `eqi` output

`ropoly-cmd <input format> eqi <path/to/original> <path/to/modified [min-gadget-length <non-negative integer>] [max-gadget-length <non-negative integer>]`

Valid input formats are `file`, `pid`, `bindump`, and `fingerprint`.

Optionally, you can specify the lengths of gadgets to be included in the fingerprint. See "Specifying gadget length" for details.

##### Gadget definition

A ROP gadget is defined as a series of consecutive instructions (consecutive both in terms of address, and in the sense
that it must be possible to consecutively execute them starting from the first instruction) ending with a return,
but for our purposes gadgets are more broadly defined to also include series ending with certain jumps and syscalls.

A gadget's length is one less than the number of included instructions--a gadget consisting of only a return is a 0-length gadget.
Usually only gadgets below a certain length are considered useful for attackers.

##### EQI calculation

Currently, EQI is calculated as the average of each of the original binary's gadgets' EQI contribution. For a gadget `g`
such that the modified binary contains no identical gadget to `g`, `g`'s EQI contribution is 0. Otherwise, `g`'s EQI
contribution is calculated as `100 * (1 - (m/t))` where `t` is the total number of gadgets in the original binary and
`m` is the size of the largest subset of gadgets from the original binary including `g` such that an offset `f` exists,
such that for each gadget `h` in the subset, the modified binary contains an identical gadget offset by `f` bytes from
its original location.

This is equivalent to the `shared-offsets` EQI function from https://github.com/polyverse/ropoly. None of the other EQI
functions found there are currently implemented.

### Specifying gadget length

When getting `fingerprint` or `eqi` output, you can add optional arguments as indicated after the required arguments.
`min-gadget-length` specifies the minimum gadget length (see "Gadget definition" under "eqi") used to generate
fingerprints and calculate EQI, and `max-gadget-length` specifies the maximum gadget length. By default, the minimum
gadget length is 0 and hte maximum gadget length is 2.

You cannot specify gadget length when using `fingperint` as the input type, because the fingerprint has already been
generated and the information needed to generate differently-sized gadgets is lost.

Attempting to get an EQI score from two fingerprints with differing minimum or maximum gadget lengths will fail.