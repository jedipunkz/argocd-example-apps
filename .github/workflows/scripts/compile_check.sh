#/bin/bash

for dirname in ./bots/*/ ; do
    cd $dirname
    if ! go build; then
        exit 1
    fi
    cd ../..
done
