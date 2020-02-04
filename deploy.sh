#!/bin/bash

set -e

read -r -p "deploy jobs? [y/n] " response
case "$response" in
    [yY][eE][sS]|[yY])
        echo "building jobs"
        CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dongfeng-jobs .

        echo "scping jobs"
        scp dongfeng-jobs dongfeng:/root/dongfeng/jobs
        scp dongfeng-jobs dongfeng-2:/root/dongfeng/jobs

        echo "cleaning jobs"
        rm dongfeng-jobs

        ;;
    *)
esac

echo "done"