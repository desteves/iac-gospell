# Pulumi program + TDD via automaiton

This Go program deploys a Pulumi program to GCP. 

To test the Go program, we deploy and destroy a Pulumi stack using the Pulumi Automation API:

- It first creates a new stack for our program.
- Then it uses Up method to deploy our stack and stream the logs to the console.
- After the deployment, it prints the URL of the website.
- Finally, it destroys the stack along with all of its the Destroy method, again streaming the logs to the console.


# Set up 

## Pre-reqs
1. Valid GCP credentials stored in `secrets/credentials.json` .
2. A valid Pulumi Cloud token stored in `secrets/pulumi_access_token.txt`

## Build
```bash
DOCKER_DEFAULT_PLATFORM=linux/amd64 DOCKER_BUILDKIT=1 docker build . --secret id=PULUMI_ACCESS_TOKEN,src=secrets/pulumi_access_token.txt  --secret id=GOOGLE_CREDENTIALS,src=secrets/credentials.json --tag nullstring/iac-gospell:latest
```


## Run 
```bash
DOCKER_DEFAULT_PLATFORM=linux/amd64 docker run nullstring/iac-gospell:latest
```

# Future work
- Have TDD use the local stack
- Add prod stack
- Add CI/CD pipeline