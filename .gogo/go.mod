module github.com/morganhein/gogo/_gogo

go 1.23.4

replace github.com/2bit-software/gogo/pkg/gogo => ./../pkg/gogo

replace github.com/2bit-software/gogo => ./../

require (
	github.com/2bit-software/gogo v0.0.0-00010101000000-000000000000
	github.com/2bit-software/gogo/pkg/gogo v0.0.0-00010101000000-000000000000
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/jessevdk/go-flags v1.6.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mvdan/sh v2.6.4+incompatible // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli/v2 v2.27.5 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/term v0.25.0 // indirect
	mvdan.cc/sh v2.6.4+incompatible // indirect
)
