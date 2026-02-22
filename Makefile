.PHONY: test wasm clean

test:
	go test -bench=. -benchmem ./...

clean:
	rm -f *.out *.wasm

dblchk.wasm:
	tinygo build -target=wasip1 -no-debug -gc=leaking -scheduler=none -panic=trap -opt=2 -o dblchk.wasm ./wasm/wasm.go

wasm: clean dblchk.wasm
