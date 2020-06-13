# build go binary
GOOS=linux go build -o main ./cmd/lambda

# create zip for aws lambda
zip api.zip main