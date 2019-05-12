#!/bin/sh

cd /root/dongfeng/sunshine/sunshinerecipenutrition

PHANTOMJS_EXECUTABLE=/usr/local/bin/phantomjs /usr/local/bin/casperjs --web-security=no --ssl-protocol=any --ignore-ssl-errors=yes main.js

cd /root/dongfeng/jobs/

./dongfeng-jobs recipe_nutrition_upload --recipe_nutrition_file_dir="/root/dongfeng/sunshine/sunshinerecipenutrition/output"

cd /root/dongfeng/sunshine/sunshinerecipenutrition/output

# clear
rm -f *.csv