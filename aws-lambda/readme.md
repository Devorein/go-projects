# Instructions 

Make sure you have aws cli installed. Then login to your aws account in the aws cli and run the following commands

**NOTE:** Replace everything between `<>` to your desired value

1. `<role-name>`: Name of the role that you want to create
2. `<function-name>`: Name of the function you want to deploy. It must match with the name of your go module
3. `<zip-file>`: Zip file that contains the `main.exe` file and will be uploaded to lambda
4. `<aws-account-id>`: Your aws account it

## Set correct env variables

```sh
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
```

## Build and package go file

```sh
go build -o main main.go
zip <zip-file> main
```

## Create a aws role first

```sh
aws iam create-role --role-name <role-name> --assume-role-policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}'
```

## Attach permission to role

```sh
aws iam attach-role-policy --role-name <role-name> --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

## Create the lambda function

```sh
aws lambda create-function --function-name <function-name> --zip-file fileb://<zip-file> --handler main --runtime go1.x --role arn:aws:iam::<aws-account-id>:role/<role-name>
```

## Invoke the lambda function

```sh
aws lambda invoke --function-name <function-name> --cli-binary-format raw-in-base64-out --payload '{"name": "John Doe", "age": 33}' output.txt
```