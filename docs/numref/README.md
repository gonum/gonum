# Numerical Reference and Worked Examples

This directory contains a Sourcey-generated numerical reference for three core
Gonum packages — **mat**, **floats**, and **stat** — along with runnable
worked examples that illustrate typical usage patterns.

## What is included

### Generated reference

The reference is built from the Go source using [Sourcey](https://sourcey.github.io/) with its Go documentation adapter:

```sh
sourcey godoc --module ./gonum \
    --packages ./mat,./floats,./stat \
    --out godoc.json

sourcey build --out .
```

It provides an extracted, browsable index of function, type, and method
entries across the three packages, with direct links back to the source
code on GitHub.

### Worked examples

Two self-contained examples demonstrate common numerical workflows:

| Example | Package | Description |
|---|---|---|
| [Solving a 2×2 linear system](examples/matrix-solve.go) | `mat` | Uses `VecDense.SolveVec` to solve a small linear system, computes the residual, and inspects the condition number. |
| [Linear regression with inverse-variance weights](examples/regression.go) | `stat` | Fits a line through measured data with known per-point uncertainties, comparing weighted and unweighted fits. |

## Running the examples

```sh
go run ./docs/numref/examples/
```

## License

All content in this directory is distributed under the BSD-style license
found in the repository root (COPYRIGHT).

## Acknowledgements

The Sourcey configuration and example material were prepared for
[Frantic bounty #33](https://gofrantic.com/bounties/p-8b91e1ac8c) and
adaptable for inclusion in Gonum's documentation.
