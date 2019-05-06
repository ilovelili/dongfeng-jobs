#!/bin/sh

cd /root/dongfeng/sunshine/sunshinemenu

PHANTOMJS_EXECUTABLE=/usr/local/bin/phantomjs /usr/local/bin/casperjs --web-security=no --ssl-protocol=any --ignore-ssl-errors=yes main.js

cd /root/dongfeng/jobs/

./dongfeng-jobs menu_upload --menu_file_path="/root/dongfeng/sunshine/sunshinemenu/output/output.csv"
