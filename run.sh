#!/bin/bash

# go build is smart enough not to include test files
# go run includes test files if we run it
go build -o bookings cmd/web/*.go && ./bookings

#if we haven't used && if the build fails then it'll run the previous version of build
# in this case second command wont be executed if the first one in not successful

# chmod +x <file_name>
# above command sets the <file_name> to be executable

# command to build the application=> ./run.sh

# command to run mailhog
# 1. download once
go get github.com/mailhog/MailHog

# 2. run this
mailhog \
  -api-bind-addr 127.0.0.1:8025 \
  -ui-bind-addr 127.0.0.1:8025 \
  -smtp-bind-addr 127.0.0.1:1025