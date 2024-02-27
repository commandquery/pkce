# Login TODO

- [ ] could definitely do logging a lot better !
- [ ] can we do this without any JavaScript?
- [ ] remove home.html; it can't see the cookies anyway
- [ ] test what happens when the JWT expires (set low JWT expiry time, eg 30 seconds)

- [X] why am I using cgo in the build process? CGO_ENABLED=1
- [X] rename repo from "login" to "pkcs"
- [X] environment variable prefixes (JWT_) seem inappropriate
- [X] make this a public repository and add a proper README
- [X] set up new gchr personal access token / delete old one
- [X] store in k8s - see deploy.yaml
- [X] deploy login to k8s
- [X] log in!
- [X] add "Secure; HTTPOnly; Domain=xxx" to cookie
- [X] store the session in a session cookie (nothing else will work)
- [X] deploy hello-world and add a JWT decoder to it
