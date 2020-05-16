# Authentication service

Auth-Service is an API based authentication service, allows authentication to your site or service with minimum efforts. 
Auth-service supports multiple authentication methods across various data sources.

## Quick Start with docker-compose

Clone auth-service repository

```
git clone https://github.com/maximthomas/blazewall.git
```

Then goto `blazewall/auth-service` directory and run

```
docker-compose up
```

This command will create three services
1. `auth-service` - authentication API service
1. `auth-service-ui` - UI client to auth-service, build with react
1. `mongo` - mongo database for user and session storage

Open http://localhost:3000 in your browser to Sign Up after signing up you can sign in with created account.

## Motivation

Blazewall auth service was created to make easier to build an authentication system.

## Authentication methods
* Username and password - authenticates exsisting user data store
* Registration - creates user in a user data stroe
* Kerberos - uses kerberos authentication

## Supported Data Sources
* LDAP
* SQL databases (in development)
* NoSQL
    * MongoDB

## Main concepts

With blazewall you can build authentication system with complexity.

There are could be different realms - for example, staff realm for employees and clients realm for the clients. All realms use their own user data stores. For example, for staff users we will use enterpse LDAP user directory, for clients we could use other database, for example, MongoDB.

There are authentication modules and authentication chains in a realm

### Authetication Module

Single authentication module, responsible for authentication or authorization.

### Authentication Chains

Authentication modues organized in authentication chains. Authentication chains is the sequence of authentication modules to organize complex authentication process.

For example, we have two modules: Login module - promts user to provide login and password and OTP module - which sends SMS with one time password to the user.

When user tries to authenticate he will be promted to enter user name and password. If the credntials are correct authentication service sends OTP via SMS and promts user to enter one time password as secon authentication factor.

In other hand, we can organize kerberos and login and password in the same chain. So if the user was not authenticatid via Kerberos he will bi prompted to enter his credentials manually.

