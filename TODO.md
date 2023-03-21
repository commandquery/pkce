# Login TODO

- [ ] add "Secure; HTTPOnly; Domain=xxx" to cookie
- [ ] store the session in a session cookie (nothing else will work)
- [ ] redirect to /apps/whoami (will print the cookie)
- [ ] deploy hello-world and add a JWT decoder to it
- [ ] work out what happens when the JWT expires (set low JWT expiry time, eg 30 seconds)
      can we do something in Go which refreshes the expired JWT? page reload is no good.

- [X] set up new gchr personal access token / delete old one
- [X] store in k8s - see deploy.yaml
- [X] deploy login to k8s
- [X] log in!