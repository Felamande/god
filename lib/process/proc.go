package process

type Process struct {
	name string
	pid  int
}

func KillByName(name string) error {
	return killByName(name)
}
