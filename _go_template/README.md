# APP Name
Description:

Maintainer team:

# APIs
OpenAPI link

# Architecture
## Layering
presenter <> internal/usecase & internal/model <> internal/infrastructure <> internal/infrastructure/store or /mq or /externalapi <> DB/PubSub/APIs

## HLD
hld

# ClickUp
ClickUp link
# How to contribute
Describe how to submit a merge request, who's the the pic, minimum number of approval, etc

# How to submit issue
Describe any convention, step to repro, any recommendation, etc


# How to use
## Run the app
- `docker-compose up -d`
- `make compile`
- `./bin/[app-name] start`

## Migrations
- `./bin/[app-name] migrate:create`
- Edit both of the generated `.sql` file
- `./bin/[app-name] migrate:run`