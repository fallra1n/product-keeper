#!/usr/bin/env bash

search_dir="./internal/core"
mock_dir="./internal/mocks"

rm -r $mock_dir/*

find "$search_dir" -type f -name "ports.go" -print0 | while IFS= read -r -d '' file; do
    dir=$(dirname "$file")
    name=$(basename $(basename $dir))
    mockgen -destination=$mock_dir/${name}/${name}.go -source=$file -package=mock${name}
done

if [ $? -ne 0 ]; then
    echo "could not find files named ports.go"
fi
