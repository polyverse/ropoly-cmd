# DEPRECATION NOTICE

Please note that this repository has been deprecated and is no longer actively maintained by Polyverse Corporation.  It may be removed in the future, but for now remains public for the benefit of any users.

Importantly, as the repository has not been maintained, it may contain unpatched security issues and other critical issues.  Use at your own risk.

While it is not maintained, we would graciously consider any pull requests in accordance with our Individual Contributor License Agreement.  https://github.com/polyverse/contributor-license-agreement

For any other issues, please feel free to contact info@polyverse.com

---

# polyverse/ropoly-cmd

## Build instructions for Docker

`build.sh` produces both a binary and a Docker image each called `ropoly-cmd`. The Docker image has the `ropoly-cmd`
binary and a directory `TestFiles` containing some simple fake fingerprints for testing the EQI calculation.

## Usage

`./ropoly-cmd <command> [flags]`

Supported commands include `fingerprint` to generate a fingerprint from a binary and output it to stdout, and
`eqi` to calculate an EQI from two fingerprints saved as files.

## Fingerprints

A fingerprint contains the gadgets (within a specified minimum and maximum length in instructions) taken from a binary.
The `fingeprint` command outputs a fingerprint as a JSON object, which if saved to a file can be used as input for the
`eqi` command.

## EQI

A number between 0 and 100 inclusive representing the difference in gadgets between a modified binary and an original binary,
with 0 being the least different and 100 meaning that the two binaries share no gadgets whatsoever.

By default, EQI is calculated as the average of each of the original binary's gadgets' EQI contribution. For a gadget `g`
such that the modified binary contains no identical gadget to `g`, `g`'s EQI contribution is 0. Otherwise, `g`'s EQI
contribution is calculated as `100 * (1 - (m/t))` where `t` is the total number of gadgets in the original binary and
`m` is the size of the largest subset of gadgets from the original binary including `g` such that an offset `k` exists,
such that for each gadget `h` in the subset, the modified binary contains an identical gadget offset by `k` bytes from
its original location.

You can change the EQI calculation to one of several using the `--eqi-func` or `-f` flag.

### `eqi-func` options

`shared-offsets` Use the default calculation.

`kill-rate` Use the percentage of gadgets from the original binary that exist at the same address in the modified binary.

`kill-rate-without-movement` Use the percentage of gadgets whose byte sequences do not appear anywhere in the modified binary's executable segments.

`highest-offset-count` Find the greatest number `n` of gadgets from the original such that an offset `k` exists and each gadget's byte sequence can be found in the modified binary at `gadget's original address`+`k`. Return 100*(1-`n`)/`total number of gadgets in original binary`.

`monte-carlo` Optionally, the flags `--trials` and `--num-gadgets` can be supplied followed by non-negative integer values.
Their defaults are 10,000 and 3 respectively. Randomly selects `--num-gadgets` gadgets from the original binary, and checks
whether an offset `k` exists such that each gadget can be found in the modified binary at its original address + `k`.
Repeats this test `--trials` times, and returns the percentage of tests in which no common offset was found.

## Gadget definition

A ROP gadget is a series of consecutive instructions (consecutive both in terms of address, and in the sense
that it must be possible to consecutively execute them starting from the first instruction) ending with a return,
but for our purposes gadgets are more broadly defined to also include series ending with certain jumps and syscalls.
A gadget is defined by both the series of instructions and the address of the first instruction.

A gadget's length is one less than the number of included instructions--a gadget consisting of only a return is a 0-length gadget.
Usually only gadgets below a certain length are considered useful for attackers.
