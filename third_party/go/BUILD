package(default_visibility = ["PUBLIC"])

go_module(
    name = "go-flags",
    module = "github.com/thought-machine/go-flags",
    version = "v1.6.0",
    deps = [
        ":unix",
    ],
)
go_module(
    name = "humanize",
    module = "github.com/dustin/go-humanize",
    version = "v1.0.0",
)

go_module(
    name = "testify",
    module = "github.com/stretchr/testify",
    install = [
        "assert",
        "require",
    ],
    version = "v1.6.0",
    test_only = True,
    deps = [
        ":difflib",
        ":spew",
        ":yaml",
    ],
)

go_module(
    name = "spew",
    module = "github.com/davecgh/go-spew",
    install = ["spew"],
    version = "v1.1.0",
    test_only = True,
)

go_module(
    name = "yaml",
    module = "gopkg.in/yaml.v3",
    version = "eeeca48fe7764f320e4870d231902bf9c1be2c08",
    test_only = True,
)

go_module(
    name = "difflib",
    module = "github.com/pmezard/go-difflib",
    install = ["difflib"],
    version = "v1.0.0",
    test_only = True,
)

go_module(
    name = "logging",
    module = "gopkg.in/op/go-logging.v1",
    version = "b2cb9fa56473e98db8caba80237377e83fe44db5",
)

go_module(
    name = "term",
    module = "golang.org/x/term",
    install = ["."],
    version = "v0.5.0",
    deps = [":unix"],
)

go_module(
    name = "unix",
    module = "golang.org/x/sys",
    install = ["unix"],
    version = "1b2967e3c290b7c545b3db0deeda16e9be4f98a2",
)
