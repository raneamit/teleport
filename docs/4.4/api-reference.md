---
title: Teleport API Reference
description: The detailed guide to Teleport API
---

# Teleport API Reference

Teleport is currently working on documenting our API.

!!! warning

        We are currently working on this project. If you have an API suggestion, [please complete our survey](https://docs.google.com/forms/d/1HPQu5Asg3lR0cu5crnLDhlvovGpFVIIbDMRvqclPhQg/edit).

## Authentication
In order to interact with the API, you will need to provision appropriate
TLS certificates. In order to provision certificates, you will need to create a
user with appropriate permissions. You should only give the api user permissions for what it actually needs. 

For example, an auditing program could use this role:

```yaml
kind: role
version: v3
metadata:
  name: auditor
spec:
  options:
    # max_session_ttl defines the TTL (time to live) of SSH certificates 
    # issued to the users with this role.
    max_session_ttl: 1h
  allow:
    logins: ['auditor']
    rules:
    - resources:
      - session
      verbs:
      - list
      - read
  deny:
    node_labels:
      '*': '*'
```

To quickly get started with the api, you can use this api-admin user.

```yaml
{!examples/go-client/api-admin.yaml!}
```

Create the user and authentication files.

```bash
$ tctl create api-admin.yaml
$ tctl auth sign --format=tls --user=api-admin --out=api-admin 
```

The above sequence should result in three PEM encoded files being generated:
`auth.crt`, `auth.key`, and `auth.cas` (certificate, private key, and CA certs respectively).

Note: by default, tctl auth sign produces certificates with a relatively short lifetime.
For production deployments, the --ttl flag can be used to ensure a more practical
certificate lifetime.

# gRPC APIs

## Client QuickStart

Follow the authentication section above to create a user and TLS certificatges. Put the generated PEM files in a folder called cert, and provide this folder in the same directory as the following go file.

```go
{!examples/go-client/client.go!}
```

## Audit Events API
Coming Soon

## Certificate Generation API
Coming Soon

## Roles APIs

[Roles](http://localhost:6600/enterprise/ssh-rbac/#roles) are used to define what resources and actions a user is allowed/denied access to.

#### Create a new role

```go
import github.com/gravitational/teleport/lib/services

// creates a new dev role which can log into teleport as 'teleport'
roleSpec := services.RoleSpecV3{
  Allow: services.RoleConditions{
    Logins: []string{"teleport"},
  },
}

role, err := services.NewRole("dev", roleSpec)
if err != nil {
  return err
}

err := client.UpsertRole(ctx, role)
```

#### Retrieve roles

```go
roles, err := client.GetRoles()
```

#### Update role

```go
role, err := client.GetRole("dev")
if err != nil {
  return err
}
  
// updates the dev role to deny access to nodes labeled as production
role.SetNodeLabels(services.Deny, services.Labels{"environment": "production"})
err = client.UpsertRole(ctx, role)
```

#### Delete role

```go
err := client.DeleteRole(ctx, "dev")
```

## Tokens API

[Tokens](http://localhost:6600/admin-guide/#adding-nodes-to-the-cluster) are used to add nodes to a cluster.

#### Create a new token

```go
tokenString, err := client.GenerateToken(ctx, auth.GenerateTokenRequest{
  // You can provide 'Token' for a non-random tokenString
  Roles: teleport.Roles{teleport.RoleProxy},
  TTL:   time.Hour,
})
```

#### Retrieve tokens

```go
tokens, err := client.GetTokens()
```

#### Update token

```go
provisionToken, err := client.GetToken(tokenString)
if err != nil {
  return err
}

// updates the token to be a proxy token
provisionToken.SetRoles(teleport.Roles{teleport.RoleProxy})
err = client.UpsertToken(provisionToken)
```

#### Delete token

```go
err := client.DeleteToken(tokenString)
```

## Workflow API
Coming Soon
