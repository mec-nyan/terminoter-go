package args

import (
	"flag"
)

type Options struct {
	File string
}

func ParseArgs() Options {
	// For now we only have this option, but there may be more in the future.  I.e. we can select a different configuration file from CLI options.
	file := flag.String("f", "", "file to read-write notes")

	flag.Parse()

	return Options{File: *file}
}
