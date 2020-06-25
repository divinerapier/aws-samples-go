#!/bin/bash

aws iam create-role \
    --role-name lambda-books-executor \
    --assume-role-policy-document file:///$(pwd)/create-lambda-role.json

## Output
#
# {
#     "Role": {
#         "RoleId": "AROA5ZNYJE74HJ3O4YYRZ",
#         "RoleName": "lambda-books-executor",
#         "Arn": "arn:aws-cn:iam::947963635704:role/lambda-books-executor",
#         "Path": "/",
#         "AssumeRolePolicyDocument": {
#             "Statement": [
#                 {
#                     "Effect": "Allow",
#                     "Principal": {
#                         "Service": "lambda.amazonaws.com"
#                     },
#                     "Action": "sts:AssumeRole"
#                 }
#             ],
#             "Version": "2012-10-17"
#         },
#         "CreateDate": "2020-06-22T02:19:58Z"
#     }
# }