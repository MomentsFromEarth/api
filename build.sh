# build go binary
GOOS=linux go build -o main

# create zip for aws lambda
zip api.zip main