subinclude("//build_defs:cuda")

cuda_binary(
    name = "cuda_main",
    srcs = ["main.cu"],
    deps = ["//:libsha256"],
)

go_binary(
    name = "go_main",
    srcs = ["main.go"],
    deps = [
        ":vso_hash",
        "//:vso-hash-cuda",
    ],
)

go_module(
    name = "vso_hash",
    module = "github.com/peterebden/vso-hash",
    version = "v1.0.0",
)
