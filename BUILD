go_library(
    name = "cli",
    srcs = [
        "flags.go",
        "logging.go",
    ],
    deps = [
        ":go-flags",
        ":humanize",
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
    get = "github.com/thought-machine/go-flags",
    revision = "v1.4.0",
)

go_get(
    name = "testify",
    get = "github.com/stretchr/testify",
    install = [
        "assert",
        "require",
    ],
    revision = "v1.6.0",
    test_only = True,
    deps = [
        ":difflib",
        ":spew",
        ":yaml",
    ],
)

go_get(
    name = "spew",
    get = "github.com/davecgh/go-spew/spew",
    revision = "v1.1.0",
    test_only = True,
)

go_get(
    name = "difflib",
    get = "github.com/pmezard/go-difflib/difflib",
    revision = "v1.0.0",
    test_only = True,
)

go_get(
    name = "humanize",
    get = "github.com/dustin/go-humanize",
    revision = "v1.0.0",
)

go_get(
    name = "yaml",
    get = "gopkg.in/yaml.v3",
    revision = "eeeca48fe7764f320e4870d231902bf9c1be2c08",
    test_only = True,
)
