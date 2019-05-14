load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "in_memory.go",
        "types.go",
    ],
    importpath = "github.com/example/project/lib/event_store",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/pkg/errors:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = ["in_memory_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/stretchr/testify/require:go_default_library",
    ],
)
