#!/bin/bash

set -e

cd /root/dongfeng/sunshine/sunshinerecipe

PHANTOMJS_EXECUTABLE=/usr/local/bin/phantomjs /usr/local/bin/casperjs --web-security=no --ssl-protocol=any --ignore-ssl-errors=yes main.js

for i in `find "$(cd /root/dongfeng/sunshine/sunshinerecipe/output; pwd)" -name "*.xls"`; do
    # https://stackoverflow.com/questions/965053/extract-filename-and-extension-in-bash
    f="${i%.*}"
    echo "converting xls to csv ... | source => $f"
    ssconvert -S $i $f.csv
done

cd /root/dongfeng/jobs/

# --recipe_file_dir=/home/min/Projects/dongfeng-sunshine-recipe-crawler/output
./dongfeng-jobs recipe_upload --recipe_file_dir="/root/dongfeng/sunshine/sunshinerecipe/output"

cd /root/dongfeng/sunshine/sunshinerecipe/output
# clear all
rm -f *.csv.* *.xls
