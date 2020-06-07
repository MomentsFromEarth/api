echo "Start Deploying to AWS Lambda"
aws lambda update-function-code --function-name mfeApi --zip fileb://api.zip --region us-east-1
echo "Done Deploying to AWS Lambda"