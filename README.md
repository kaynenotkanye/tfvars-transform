This is a Golang command line utility which gives the ability to transform/assign values to a key/value file, such as a terraform.tfvars file. The initial scope of this utility was specifically made to find and replace AMI IDs when given
a specific key name. A modified regex can broaden the scope and allow this tool to be more versatile.

Example terraform.tfvars file:
```
aws_region = "us-west-2"
aws_account_id = "redacted-abcd1234"
appname_ami_id = "ami-abcd1234"
```

Example usage:
```
tfvars-transform -f terraform.tfvars -k appname_ami_id -v ami-mynewamid
```

Resulting terraform.tfvars (Notice that the utility will only update the value of the specified key passed through the -k flag):
```
aws_region = "us-west-2"
aws_account_id = "redacted-abcd1234"
appname_ami_id = "ami-mynewamid"
```

Another example command (If you have an ami-id stored in a text file):
```
tfvars-transform -f terraform.tfvars -k appname_ami_id -v $(cat ami-id.txt)
```
