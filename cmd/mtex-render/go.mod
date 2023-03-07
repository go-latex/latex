module github.com/go-latex/latex/cmd/mtex-render

go 1.19

require (
	gioui.org v0.0.0-20210822154628-43a7030f6e0b
	github.com/go-fonts/latin-modern v0.3.0
	github.com/go-fonts/liberation v0.3.0
	github.com/go-latex/latex v0.0.0-20230307184137-b3ecf8b1eeee
	golang.org/x/image v0.6.0
)

require (
	gioui.org/cpu v0.0.0-20210817075930-8d6a761490d2 // indirect
	gioui.org/shader v1.0.0 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/exp v0.0.0-20210722180016-6781d3edade3 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.8.0 // indirect
)

replace github.com/go-latex/latex => ../..
