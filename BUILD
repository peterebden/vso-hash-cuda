subinclude("//build_defs:cuda")

cuda_library(
    name = "libsha256",
    srcs = ["sha256.cu"],
    hdrs = ["sha256.h"],
    visibility = ["PUBLIC"],
)

cgo_library(
    name = "vso-hash-cuda",
    srcs = ["vsohash.go"],
    hdrs = ["sha256.h"],
    ldflags = ["-lsha256 -L. -lcudart_static -ldl -lrt -lstdc++"],
    deps = [
        ":libsha256",
    ],
    visibility = ["PUBLIC"],
)

go_test(
    name = "vsohashcuda_test",
    srcs = ["vsohashcuda_test.go"],
    external = True,
    deps = [
        ":testify",
        ":vso-hash-cuda",
    ],
)

go_benchmark(
    name = "vsohashcuda_benchmark",
    srcs = ["vsohashcuda_benchmark_test.go"],
    external = True,
    deps = [
        ":vso-hash-cuda",
    ],
)

go_module(
    name = "testify",
    install = [
        "assert",
        "require",
    ],
    licences = ["MIT"],
    module = "github.com/stretchr/testify",
    version = "v1.7.0",
    deps = [
        ":difflib",
        ":spew",
        ":yaml",
    ],
)

go_module(
    name = "difflib",
    install = ["difflib"],
    licences = ["BSD-3-Clause"],
    module = "github.com/pmezard/go-difflib",
    version = "v1.0.0",
)

go_module(
    name = "spew",
    install = ["spew"],
    licences = ["ISC"],
    module = "github.com/davecgh/go-spew",
    version = "v1.1.1",
)

go_module(
    name = "yaml",
    licences = ["MIT"],
    module = "gopkg.in/yaml.v3",
    version = "v3.0.0-20210107192922-496545a6307b",
)
