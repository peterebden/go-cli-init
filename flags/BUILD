go_library(
    name = "flags",
    srcs = ["flags.go"],
    deps = [
        "//third_party/go:go-flags",
        "//third_party/go:humanize",
    ],
)

go_test(
    name = "flags_test",
    srcs = ["flags_test.go"],
    deps = [
        ":flags",
        "//third_party/go:testify",
    ],
)