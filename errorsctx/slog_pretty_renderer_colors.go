package errorsctx

type prettyWriterColorProfile struct {
	reset    string
	bold     string
	time     string
	trace    string
	debug    string
	info     string
	warn     string
	error    string
	panic    string
	levelu   string // Инвертированная плашка для критических ошибок
	loc      string
	link     string
	stdots   string
	sttext   string
	key      string
	errkey   string
	errmeta  string // Тёмно-оранжевый (TrueColor) для метаданных ошибок
	errstage string // Оранжевый (TrueColor) для стадий/слоев ошибок
	ctx      string
}

func newPrettyWriterColorProfileDark() *prettyWriterColorProfile {
	return &prettyWriterColorProfile{
		reset:  "\x1b[0m",
		bold:   "\x1b[1m",
		time:   "\x1b[35m",
		trace:  "\x1b[90m",
		debug:  "\x1b[36m",
		info:   "\x1b[32m",
		warn:   "\x1b[33m",
		error:  "\x1b[31m",
		panic:  "\x1b[1;41;97m",
		levelu: "\x1b[1;41;97m",
		loc:    "\x1b[38;5;244m",

		// КОРРЕКЦИЯ: Поднимаем яркость палочек и двоеточий до читаемого темно-серого (240)
		link:   "\x1b[38;5;240m",
		stdots: "\x1b[38;5;240m",

		sttext:   "\x1b[38;5;245m",
		key:      "\x1b[38;5;109m",
		errkey:   "\x1b[38;5;203m",
		errmeta:  "\x1b[38;2;255;140;0m",
		errstage: "\x1b[38;2;255;165;0m",
		ctx:      "\x1b[38;5;252m",
	}
}

func newPrettyWriterColorProfileLight() *prettyWriterColorProfile {
	return &prettyWriterColorProfile{
		reset:  "\x1b[0m",
		bold:   "\x1b[1m",
		time:   "\x1b[95m",
		trace:  "\x1b[90m",
		debug:  "\x1b[36m",
		info:   "\x1b[32m",
		warn:   "\x1b[33m",
		error:  "\x1b[31m",
		panic:  "\x1b[1;41;97m",
		levelu: "\x1b[1;41;97m",
		loc:    "\x1b[38;5;240m",

		// КОРРЕКЦИЯ: Для светлой темы делаем разделители более темными и контрастными
		link:   "\x1b[38;5;244m",
		stdots: "\x1b[38;5;244m",

		sttext:   "\x1b[38;5;240m",
		key:      "\x1b[38;5;31m",
		errkey:   "\x1b[38;5;203m",
		errmeta:  "\x1b[38;2;255;140;0m",
		errstage: "\x1b[38;2;255;165;0m",
		ctx:      "\x1b[38;5;238m",
	}
}
