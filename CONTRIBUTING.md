# Contributing

## Issues

Log an issue for any question or problem you might have. When in doubt, log an issue, and
any additional policies about what to include will be provided in the responses. The only
exception is security disclosures which should be sent privately.

Committers may direct you to another repository, ask for additional clarifications, and
add appropriate metadata before the issue is addressed.

## Contributions

Any change to resources in this repository must be through pull requests. This applies to all changes
to documentation, code, binary files, etc.

No pull request can be merged without being reviewed and approved.

If you are looking to contribute to zoau Node.js development, follow these steps
to set up your development environment:

### Cloning and Building 

1. Follow the instructions in
https://www.ibm.com/docs/en/zoau/latest?topic=installing-configuring-zoa-utilities to install
and configure ZOAU on your system.

After installation, make sure the PATH and LIBPATH environment variables include the location
of the ZOAU binaries and dlls as follows:
``` bash
export PATH=<path_to_zoau>/bin:$PATH
export LIBPATH=<path_to_zoau>/lib:$LIBPATH
```

2. Clone the zoau-node repository.

```bash
$ git clone git@github.com/Stolkerve/zoau-go
```

3. Install the dependencies required for zoau Node.js development.

```bash
$ cd zoau-go
```

## Tests

Verify that the zoau Node.js module s working by running the test suite.

```bash
$ go test -v
```

### Commit message

A good commit message should describe what changed and why.

It should:
  * contain a short description of the change
  * be entirely in lowercase with the exception of proper nouns, acronyms, and the words that refer to code, like function/variable names
  * be prefixed with one of the following words:
    * fix: bug fix
    * hotfix: urgent bug fix
    * feat: new or updated feature
    * docs: documentation updates
    * refactor: code refactoring (no functional change)
    * perf: performance improvement
    * test: tests and CI updates