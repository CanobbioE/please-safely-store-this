rm -r ./build/*
for GOOS in darwin linux windows; do
    for GOARCH in 386 amd64; do
        go build -v -o "build/psst-$1-$GOOS-$GOARCH"
    done
done
