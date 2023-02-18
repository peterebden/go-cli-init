go_library(
    name = "logging",
    srcs = [
        "logging.go",
    ],
    deps = [
        "//third_party/go:logging",
        "//third_party/go:term",
    ],
)

go_test(
    name = "logging_test",
    srcs = ["logging_test.go"],
    deps = [
        ":logging",
        "//third_party/go:logging",
        "//third_party/go:testify",
    ],
)
