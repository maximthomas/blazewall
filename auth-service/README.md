# Authentication service

Auth-Service is an API based authentication service, allows adding authentication to your site or service with minimum efforts. 
Auth-service supports multiple authentication methods across various data sources.
It allows building complex authentication process with various steps or different authentication methods.   

## Quick Start with docker-compose

Clone auth-service repository

```
git clone https://github.com/maximthomas/blazewall.git
```

Then goto `blazewall/auth-service` directory and run `docker-compose`

```
docker-compose up
```

This command will create three services
1. `auth-service` - authentication API service
1. `auth-service-ui` - UI client for auth-service, built with React
1. `mongo` - Mongo database for user and session storage

Open http://localhost:3000 in your browser to Sign Up. After signing up you can Sign In with created account.

## Motivation

Blazewall auth service was created to make easier to build an authentication system.
<!-- TODO extend motivation-->

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

With Blazewall you can build authentication system with any desired complexity.

There are could be different realms - for example, staff realm for employees and clients realm for the clients. All realms use their own user data stores. For example, for staff users we will use enterpse LDAP user directory, for clients we could use other database, for example, MongoDB.

There are authentication modules and authentication chains in a realm

### Authetication Module

Single authentication module, responsible for authentication or authorization step.
For example - prompt username and password or send and verify one-time password.

### Authentication Chains

Authentication modules organized in authentication chains. 
Authentication chain is the sequence of authentication modules to organize complex authentication process.
For example, we have two modules: Login module - prompts user to provide a login and password and OTP module - which sends SMS with one time password to the user.

When user tries to authenticate he will be prompted to enter user name and password. 
If the credentials are correct authentication service sends OTP via SMS and prompts user to enter one time password as a second authentication factor.
In other hand, we can organize kerberos and login and password in the same chain. 
So if the user was not authenticated via Kerberos he will bi prompted to enter his credentials manually.

## Configuration Reference

```yaml
authentication: #section defines everything related to authentication process 
  realms: # defines real
    users: #realm ID
      modules: # authentication modules
        login: # authentication module ID - used in authentication chain
          type: "login" # could be "login", "registration", "kerberos"
          properties: #module properties map
        registration:
          type: "registration"
          properties:
            additionalFileds:
              - dataStore: "name"
                prompt: "Name"

      authChains: # defines authentication chains
        login: # authentication chain ID
          modules: # authentication chain modules list
            - id: "login" # module id
              properties:
        registration:
          modules:
            - id: "registration"
              properties:

      userDataStore: # defines User Data Stokre
        type: "mongodb" # could be "mongodb" or "ldap"
        properties:
          url:  "mongodb://root:changeme@localhost:27017"
          database:   "users"
          collection: "users"
          userAttributes: # additional user attributes 
            - "name"

session:
  type: "stateless" # could be also "stateful"
  expires: 60000 #token lifetime in seconds
  jwt: #JWT properties
    issuer: 'http://auth-service'
    privateKeyPem: |
      -----BEGIN RSA PRIVATE KEY-----
      MIIBOQIBAAJATmLeD2qa5ejVKJ3rwcSJaZAeRw4CVrUHvi1uVvBah6+6qCdjvH8N
      RT+GOI3ymdnilILPHcn51A0XQAXyrvFkgwIDAQABAkAPZUvIK2ARGBIF0D6l6Dw1
      B6Fqw02iShwjNjkdykd9rsZ+UwsYHJ9xXSa2xp7eGurIUqyaDxF+53xpE9AH72PB
      AiEAlEOIScKvyIqp3ZAxjYUd3feke2AGq4ckoq/dXFvxKHcCIQCHWH+6xKyXqaDL
      bG5rq18VQR2Nj7VknY4Eir6Z6LrzVQIgSz3WbXBi2wgb2ngx3ZsfpCToEUCTQftM
      iU9srFFwmlMCIFPUbMixqHUHi6BzuLDXpDz15+gWarO3Io+NoCCUFbdBAiEAinVf
      Lnb+YDP3L5ZzSNF92P9yBQaopFCifjrUqSS85uw=
      -----END RSA PRIVATE KEY-----

  dataStore: # session data store
    type: "mongo" 
    properties:
      url: "mongodb://root:changeme@localhost:27017"
      database:   "session"
      collection: "session"
```