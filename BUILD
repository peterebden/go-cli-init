go_library(
    name = "cli",
    srcs = [
        "flags.go",
        "logging.go",
    ],
    deps = [
        ":go-flags",
        ":logging",
        ":terminal",
    ],
)

go_test(
    name = "logging_test",
    srcs = ["logging_test.go"],
    deps = [
        ":cli",
        ":logging",
        ":testify",
    ],
)

go_test(
    name = "flags_test",
    srcs = ["flags_test.go"],
    deps = [
        ":cli",
        ":testify",
    ],
)

go_get(
    name = "logging",
    get = "gopkg.in/op/go-logging.v1",
    revision = "b2cb9fa56473e98db8caba80237377e83fe44db5",
)

go_get(
    name = "terminal",
    get = "golang.org/x/crypto/ssh/terminal",
    revision = "505ab145d0a99da450461ae2c1a9f6cd10d1f447",
    deps = [":unix"],
)

go_get(
    name = "unix",
    get = "golang.org/x/sys/unix",
    revision = "1b2967e3c290b7c545b3db0deeda16e9be4f98a2",
)

go_get(
    name = "go-flags",
    get = "github.com/jessevdk/go-flags",
    revision = "v1.4.0",
)

go_get(
    name = "testify",
    get = "github.com/stretchr/testify",
    install = [
        "assert",
        "require",
        "vendor/github.com/davecgh/go-spew/spew",
        "vendor/github.com/pmezard/go-difflib/difflib",
    ],
    revision = "f390dcf405f7b83c997eac1b06768bb9f44dec18",
    deps = [":spew"],
)

go_get(
    name = "spew",
    get = "github.com/davecgh/go-spew/spew",
    revision = "ecdeabc65495df2dec95d7c4a4c3e021903035e5",
)
