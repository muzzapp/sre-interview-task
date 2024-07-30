# Deployment status service

This service is an API allowing for things that are managing deployments to write data and things representing the data to query it.

Data is written from multiple sources (e.g. Lambdas, ECS, EKS, EC2) so this exposes a uniform API for these sources to store the data.

# Setting up the development environment

## Requirements

You'll need the following tools

* Terraform
* Docker
* Go

To install these on OSX run `brew install hashicorp/tap/terraform golang docker`

## Setting up the database

Running `make dev_env` will setup a local DynamoDB and configure it with Terraform

## Generating data (optional)

Running `make data_generation` will generate some data into the local DynamoDB

## Running the API

Running `make run_all_docker_compose` will run the API in Docker, however it can also be run locally with `make run_api_local`

## Resetting the data

Running `make reset_dev_env` will reset all the data in DynamoDB

## Other tools that may be helpful

* NoSQL Workbench for viewing data, you can follow [these instructions](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/workbench.querybuilder.connect.html) to get set up. Use `localhost:4566` for the host, the region and credentials can be any value as this is using localstack.

# Using the API

## Terminology

* Component is a unique name for a service within an environment, it is always lower case alpha-numeric characters that may contain hyphens and full stops. e.g. `abc.def`, `ab-123`, and `123-def.xyz`
* Environment is a string that matches `^[a-z]+(-[a-z]+){0,2}$`
* Deployable is a unique combination of component & environment

## /records

`curl -d '{"Records":[{"Component":"abc.def","Environment":"prod-all","State":"Ok","Timestamp":"2024-07-30T14:06:41+00:00"}]}' localhost:8080/records `

## /historical/:component

`curl localhost:8080/historical/abc.def | jq`

```
{
  "Records": [
    {
      "Component": "abc.def",
      "Environment": "prod-all",
      "Timestamp": "2024-07-30T15:06:41+01:00",
      "State": "Ok"
    }
  ]
}
```

## /current

This endpoint returns the latest status for each deployable

`curl localhost:8080/current | jq`

```
{
  "Records": [
    {
      "Component": "abc.def",
      "Environment": "prod-all",
      "Timestamp": "2024-07-30T15:06:41+01:00",
      "State": "Ok"
    }
  ]
}
```

# Task summary

Imagine you are working on an engineering team and someone has put this code up for a PR for you to review.

Your task is the following;

1. Fork this repository in GitHub
2. Create a PR with comments detailing
  1. The mistakes that you find
  2. Corrections you'd suggest or a rough description of what the change should be
  3. The reasoning behind the change
2. The things that can't be changed are
  1. The body returned in normal operation by the external API defined by the API codebase must not change
  2. The data must be stored in DynamoDB, however the underlying format may be changed

If something isn't clear please do let us know so we can help clarify any issues.
