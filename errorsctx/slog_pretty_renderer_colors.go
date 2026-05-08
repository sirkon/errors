package errorsctx

type prettyWriterColorProfile struct {
	reset  string
	bold   string
	time   string
	trace  string
	debug  string
	info   string
	warn   string
	error  string
	panic  string
	loc    string // Color of locations for log itself and @location of error
	link   string // Hierarchy links of tree.
	stdots string
	sttext string
	key    string
	errkey string
	ctx    string
}

func newPrettyWriterColorProfileDark() *prettyWriterColorProfile {
	return &prettyWriterColorProfile{
		reset:  "\033[0m",
		bold:   "\033[1m",
		time:   "\033[35m",
		trace:  "\033[90m",
		debug:  "\033[36m",
		info:   "\033[32m",
		warn:   "\033[33m",
		error:  "\033[31m",
		panic:  "\033[1;41;97m",
		loc:    "\033[38;5;244m",
		link:   "\033[38;5;240m",
		stdots: "\033[38;5;236m",
		sttext: "\033[38;5;245m",
		key:    "\033[38;5;109m",
		errkey: "\033[38;5;203m",
		ctx:    "\033[38;5;252m",
	}
}

func newPrettyWriterColorProfileLight() *prettyWriterColorProfile {
	return &prettyWriterColorProfile{
		reset:  "\033[0m",
		bold:   "\033[1m",
		time:   "\033[95m",
		trace:  "\033[90m",
		debug:  "\033[36m",
		info:   "\033[32m",
		warn:   "\033[33m",
		error:  "\033[31m",
		panic:  "\033[1;41;97m",
		loc:    "\033[38;5;240m",
		link:   "\033[38;5;248m",
		stdots: "\033[38;5;252m",
		sttext: "\033[38;5;240m",
		key:    "\033[38;5;31m",
		errkey: "\033[38;5;203m",
		ctx:    "\033[38;5;238m",
	}
}
