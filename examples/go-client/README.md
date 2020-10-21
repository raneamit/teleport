## Teleport Auth Go Client

### Introduction

Teleport Auth Server has an API which hasn't been officially documented (yet).
Both `tctl` and `tsh` use the Auth API to:

* Request certificates (`tsh login` or `tctl auth sign`)
* Add nodes and users (`tctl users` and `tctl nodes`)
* Manipulate cluster state (`tctl` resources)

### API Authentication

Auth API clients must perform two-way authentication using x509 certificates:

1. They have to validate the auth server x509 certificate to make sure the
   API endpoint can be trusted.
2. They must offer their x509 certificate, which has been previously issued
   by the auth sever.

### Demo

This little program demonstrates how to:

1. Authenticate against the Auth API using two certificates.
2. Make API calls to issue CRUD requests for cluster join tokens, roles, and labels.
3. Receive, allow, and deny access requests.

First, get your teleport server running. If you aren't using the default 
teleport data directory,  change `dataDir` in main.go to match yours.

Next, you have to use `tctl` to issue an API certificate, i.e. on the auth server:

```
$ tctl auth export --type=tls > /var/lib/teleport/ca.cert
```

Now you can run the go example.

```bash
$ go get github.com/gravitational/teleport/lib/auth
$ go run .
```

This should work as long as you execute it on the same auth server:

### TODO

This Auth server API allows clients to "jump" to API endpoints of all trusted
clusters connected to it. We need to add a snippet how to enumerate trusted
clusters and connect to their API endpoints later.
