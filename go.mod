module example

go 1.22

require (
	github.com/google/uuid v1.6.0
	my-zinx v0.0.0-00010101000000-000000000000
)

require (
	github.com/jfeliu007/goplantuml v1.6.2 // indirect
	github.com/petermattis/goid v0.0.0-20250319124200-ccd6737f222a // indirect
	github.com/stretchr/testify v1.10.0 // indirect
)

replace my-zinx => ./zinx
