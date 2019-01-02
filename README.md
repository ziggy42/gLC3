# gLC3

LC-3 VM written in Go.
Inspired by [Write your Own Virtual Machine](https://justinmeiners.github.io/lc3-vm/index.html).

## Build
```bash
export GOPATH=$(pwd)
go build github.com/ziggy42/LC3
```

## Run
```bash
./LC3 /path/to/object
```

## Compile a program
```bash
docker build -t lc3compiler .
docker run --rm -ti -v $(pwd):/data lc3compiler /data/programs/hello-world.asm
```
