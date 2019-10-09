<!--
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
-->

# BLAZEWALL

![Blazewall logo](docs/img/logo.png)

**BLAZEWALL** is an Open Source Single Sign On and Access Management platform built in microservice architecture 
and released under Apache 2.0 license.

### Solution architecture shown on the diagram below:
![Services interaction diagram](docs/img/services-diagram.png)

### Services:

|Service|Description|
|-------|-----------|
|**auth-service**|Authentication frontend, responsible for signing up or signing in users.|
|**gateway-service**|Proxies all user requests to protected resources. Gateway insures, if the user request does not violate security policy and pass the request to the protected resource. If the request violates policy, gateway service denies request, and redirects user to authentication.|
|**session-service**|Stores and manages user sessions.|
|**user-service**|Responsible for user account management.|
|**policy-service**|Will be developed in future releases, to externalize policy decision to external service from gateway-service|
|**test-service**|Test service with unsecured and secured zone|


## Processes

### Authentication process diagram

![Access protected resource process](docs/img/authentication.png)

### Access Protected Resource

![Authentication process](docs/img/access-protected-resource.png)


## Quick Start

Quick start with docker-compose

Add extry to file `/etc/hosts` in Linux or Mac or `c:\Windows\System32\Drivers\etc\hosts`
```
127.0.0.1 example.com auth.example.com
```

Start all services locally with docker-compose
```bash
$ docker-compose up --build
```
After all services started go to `http://example.com:8080/`, you'll see page avialable to all users
Click on the `Try to Authenticate` button. You will be redirected to the page `http://example.com:8080/user` protected by the `gateway-servce`

`gateway-service` checks if user is authenticated and if he is not, redirects to the `auth-service` `http://auth.example.com:8081/auth-service/v1/users`

Enter login `admin` and password `password` to authenticate.

After authentication, you will be redirected back to the protected resource `http://example.com:8080/user`


