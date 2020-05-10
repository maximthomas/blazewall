# Authentication service

Auth-Service is an API based authentication service, allows authentication to your site or service with minimum efforts. 
Auth-service supports multiple authentication methods across various data sources.

## Quick Start

Create configuration file

```yaml
authentication:
  realms:
    users:
      modules:
        login:
          type: "login"
          properties:
        registration:
          type: "registration"
          properties:
            additionalFileds:
              - dataStore: "name"
                propmt: "Name"

      authChains:
        login:
          modules:
            - id: "login"
              properties:
        registration:
          modules:
            - id: "registration"
              properties:

      userDataStore:
        type: "mongodb"
        properties:
          url:  "mongodb://root:example@localhost:27017"
          database:   "users"
          collection: "users"
          userAttributes:
            - "name"

      session:
        type: "stateless"
        expires: 60000
        jwt:
          issuer: 'https://auth-service'
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

```
 
## Authentication methods
* Username and password
* Registration
* Kerberos

## Supported Data Sources
* LDAP
* SQL databases
    * PostgreSQL
    * MySQL
* NoSQL
    * MongoDB
    




 