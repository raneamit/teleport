---
title: Teleport API Reference
description: The detailed guide to Teleport API
---

# Teleport API Reference

Teleport is currently working on documenting our API.

!!! warning

        We are currently working on this project. If you have an API suggestion, [please complete our survey](https://docs.google.com/forms/d/1HPQu5Asg3lR0cu5crnLDhlvovGpFVIIbDMRvqclPhQg/edit).

## Authentication
In order to interact with the Access Request API, you will need to provision appropriate
TLS certificates. In order to provision certificates, you will need to create a
user with appropriate permissions:

```bash
$ cat > rscs.yaml <<EOF
kind: user
metadata:
  name: access-plugin
spec:
  roles: ['access-plugin']
version: v2
---
kind: role
metadata:
  name: access-plugin
spec:
  allow:
    rules:
      - resources: ['access_request']
        verbs: ['list','read','update']
    # teleport currently refuses to issue certs for a user with 0 logins,
    # this restriction may be lifted in future versions.
    logins: ['access-plugin']
version: v3
EOF
# ...
$ tctl create rscs.yaml
# ...
$ tctl auth sign --format=tls --user=access-plugin --out=auth
# ...
```

The above sequence should result in three PEM encoded files being generated:
`auth.crt`, `auth.key`, and `auth.cas` (certificate, private key, and CA certs respectively).

Note: by default, tctl auth sign produces certificates with a relatively short lifetime.
For production deployments, the --ttl flag can be used to ensure a more practical
certificate lifetime.

# gRPC APIs

## Client QuickStart

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

role, _ := services.NewRole(roleName, services.RoleSpecV3{
	Allow: services.RoleConditions{
    Logins: []string{login},
	},
})

client.UpsertRole(ctx, role)
```

#### retrieve roles

```go
roles, err := client.GetRoles()
if err != nil {
  return err
}
```

#### update role

```go
// update a role to deny access to nodes labeled as production
// you can update any attributes a role has with its setter methods
role, err = client.GetRole("dev")
if err != nil {
  log.Fatal(err)
}
  
role.SetNodeLabels(services.Deny, services.Labels{"environment": "production"})
if err = client.UpsertRole(ctx, role); err != nil {
  return err
}
```

#### delete role

```go
if err = client.DeleteRole(ctx, "dev"); err != nil {
  return err
}
```

## Tokens API

[Tokens](http://localhost:6600/admin-guide/#adding-nodes-to-the-cluster) are used to add nodes to a cluster.

#### Create a new token

```go
tokenName := "mytoken"
token, err := client.GenerateToken(ctx, auth.GenerateTokenRequest{
  Token: tokenName,
  Roles: teleport.Roles{teleport.RoleProxy},
  TTL:   time.Hour,
})
if err != nil {
  log.Fatalf("Failed to generate token: %v", err)
}
log.Printf("Generated token: %v\n", token)

// generate a random cluster join token for adding a node to a cluster
randToken, err := client.GenerateToken(ctx, auth.GenerateTokenRequest{
  Roles: teleport.Roles{teleport.RoleNode},
  TTL:   time.Hour,
})
if err != nil {
  log.Fatalf("Failed to generate random token: %v", err)
}
log.Printf("Generated random token: %v\n", randToken)
```

#### retrieve tokens

```go
tokens, err := client.GetTokens()
if err != nil {
  return err
}
```

#### update token

```go
// update a token to be a proxy token
// you can update any attributes a token has with its setter methods
provToken, err := client.GetToken(token)
if err != nil {
  log.Fatal(err)
}

provToken.SetRoles(teleport.Roles{teleport.RoleProxy})
err = client.UpsertToken(provToken)
if err != nil {
  log.Fatal(err)
}
```

#### delete token

```go
if err = client.DeleteToken(token); err != nil {
  return err
}
```

## Workflow API
Coming Soon