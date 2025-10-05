module example

go 1.22

require (
	github.com/Meha555/go-tinylog v1.0.2
	github.com/Meha555/pulse v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
)

require github.com/petermattis/goid v0.0.0-20250904145737-900bdf8bb490 // indirect

replace github.com/Meha555/pulse => ..
