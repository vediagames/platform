# How to start

This guide assumes you are on a UNIX based environment. Also make sure you're in the root of project.

## Deploy postgres

```
# make
make postgres

# no make
docker run -e "POSTGRES_USER=vedia" -e "POSTGRES_PASSWORD=123" -e "POSTGRES_DB=vediagames" -p 5432:5432 -d postgres:15.0-alpine
```

## Deploy redis

```
# make
make redis

# no make
docker run -p 6379:6379 -d redis
```

## Run gqlgen

```
# make
make gqlgen/bff

# no make
cd bff && go get github.com/99designs/gqlgen && go run github.com/99designs/gqlgen generate
```

## Create config.yml

You can pretty much copy the content from config.example.yml to config.yml

```
cp config.example.yml config.yml
```

## Run migrations

`go run . migrate`

## Run stub - this adds stub data to your development DB

`go run . stub`

## Run server

`go run . server`

# Code architecture/patterns

Quick explanations regarding some decisions and patterns established.

## Single-binary approach

We have a single main.go, which is made with cobra framework. This ensures that the whole app is pretty much a single
binary that can be created and deployed as a single container image. This makes things very manageable.

There are some helper modules that are also used, but for the bigger part, everything is pretty much here build and run
as a single binary.

## Package architecture and their meaning

I'll explain what each package contains and what is it's meaning of existence.

### bff

Short for backend-for-frontend. Basically, a gateway that exposes backend services through GraphQL API. Uses gqlgen for
GraphQL code generation. Also code for retrieving data for specific pages on vediagames.com lives here. (Let's call it
page service)

### bucket

S3 and GCS (Google Cloud Storage) client implementations for providing S3 storage functionality to vediagames. We use it
for uploading images and videos when publishing new games.

### category, game, section, tag

Services.

### cmd

Cobra commands (migrate, server, stub)

### config

Config package.

### fetcher

Client implementations for game distribution and game monetize. These provide games for us and these implementations
fetch them.

### db

Schema and migrations for the database. Also includes some stub data for stubbing local DB for development.

### search

Search service that provides search functionality for the website.
This service connects all other services that support search and adds business logic on top of them.

# Go get for private repositories

## Generate your SSH key with the `ssh-keygen` tool

<https://git-scm.com/book/en/v2/Git-on-the-Server-Generating-Your-SSH-Public-Key>

A little advice here, use different keys for different services/accounts. Prefix them with the service/account name. For
example `id_rsa.github`

## Add the key to your GitHub account

https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account

## Edit your .gitconfig

Simply run this command:

`git config --global --add url."ssh://git@github.com/".insteadOf "https://github.com/"`

### Or if you wish to edit your .gitconfig file directly:

Append this two lines to the end of the document:

```
[url "ssh://git@github.com/"]
    insteadOf = https://github.com/
```

## Edit your \~/.ssh/config

Add these 5 lines:

```
Host github.com
    User git
    Hostname github.com
    IdentitiesOnly yes
    IdentityFile ~/.ssh/id_rsa_github `**<- MAKE SURE THAT THIS POINTS TO YOUR SSH KEY THAT YOU GENERATED IN STEP #1**
```

---

This should now allow you to freely clone repos from github.com/vediagames

You can check whether everything went well by running this command (this clones vediagames.com/category repo)

`git clone git@github.com:vediagames/vediagames.com.git`

## go get

After you've done setting up your SSH authentication which allows you to freely clone repos from GH, you still need to
do something else to make everything run smoothly on your development machine (as far as `Go` development goes)

Luckily this is very simple, you can just run:

`go env -w GOPRIVATE="github.com/vediagames"`

To test, you can try running this inside any folder that has go.mod file:

`go get github.com/vediagames/zeroerror`

You can also use GoLand to set it up by editing the Environment input box in `File/Settings/Go Modules.` But that will
only work in the directory you're currently in.