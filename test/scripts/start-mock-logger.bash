#!/bin/bash
# 2018-12-27 (cc) <paul4hough@gmail.com>
#
cd $HOME/demo
echo `date` starting logger >> log/start-mock-logger.log
nohup bin/mock-logger \
  --laddr ":5002" \
  --log-fn "log/mock-logger.log" \
  > log/mock-logger.out 2>&1 &
echo `date` pid $! >> log/start-mock-logger.log
