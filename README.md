# gLC3

## Compile a program
```bash
docker build -t lc3compiler .
docker run --rm -ti -v $(pwd):/data lc3compiler /data/programs/hello-world.asm
```
