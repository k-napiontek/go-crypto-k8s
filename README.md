
S3
terraform apply -var-file=../global.tfvars   

terraform init -migrate-state backend-config=../state.config 

IAM
terraform init -backend-config=../state.config 

terraform apply -var-file=../global.tfvars   

ECR
terraform init -backend-config=../state.config 

terraform apply -var-file=../global.tfvars  