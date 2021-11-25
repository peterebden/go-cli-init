go_library(
    name = "flags",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    deps = [
        "//third_party/go:go-flags",
        "//third_party/go:humanize",
    ],
)

go_test(
    name = "flags_test",
    srcs = glob(["*_test.go"]),
    data = ["test_data"],
    deps = [
        ":flags",
        "//third_party/go:testify",
    ],
)
