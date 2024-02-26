# Login TODO

- [ ] why am I using cgo in the build process? CGO_ENABLED=1
- [ ] rename repo from "login" to "pkcs"
- [ ] environment variable prefixes (PKCE_) seem inappropriate
- [ ] remove the coachcentric-specific deploy.yaml (move it to CC)
- [ ] home.html is a useful debugging trick, document in README but make it disable-able 
- [ ] make this a public repository and add a proper README
- [ ] work out what happens when the JWT expires (set low JWT expiry time, eg 30 seconds)
      can we do something in Go which refreshes the expired JWT? page reload is no good.

- [X] set up new gchr personal access token / delete old one
- [X] store in k8s - see deploy.yaml
- [X] deploy login to k8s
- [X] log in!
- [X] add "Secure; HTTPOnly; Domain=xxx" to cookie
- [X] store the session in a session cookie (nothing else will work)
- [X] deploy hello-world and add a JWT decoder to it
