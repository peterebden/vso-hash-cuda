vso-hash-cuda
=============

Golang implementation of the BuildXL paged hash function (a.k.a. VSO-Hash)

See https://github.com/microsoft/BuildXL/blob/master/Documentation/Specs/PagedHash.md for
more information about the function in general.

See https://github.com/peterebden/vso-hash for a CPU-based implementation.


Building
--------

I (unsurprisingly) recommend building this with [Please](https://please.build).
It should be as simple as `./pleasew run //cmd:go_main` to run the Go version of it.
You will obviously have to have the CUDA tools installed (`nvidia-cuda-toolkit` etc) for
anything to work.

If you want to build it using `go build`, you will have to manually invoke `nvcc` first:
```
nvcc --shared --compiler-options='-fPIC' -o libsha256.so sha256.cu
go build
```
Unfortunately `go test` doesn't work. It's also unlikely to ever be `go get`-able since
there's no way of scripting that `nvcc` invocation. See the link above for more details
on plz which can do that!


Notes
-----

Currently this is (sadly) losing to the CPU version in terms of performance, although does have the
advantage that most of your CPU cores are still available for doing other things.
Most likely it is limited by something else at this point, e.g. CPU, memory bandwidth or general
implementation tomfoolery.
Since it's mostly a learning project at the moment I'm not highly motivated to dig into the details
of exactly what is going wrong.

The `Sequential0` test is not working; no doubt there is some silly edge case that isn't being
handled correctly (I've already fixed one of these but when it didn't immediately fix the problem
I gave up).

The contents of `sha256.cu` are heavily based on https://github.com/B-Con/crypto-algorithms/blob/master/sha256.c
by Brad Conte.
