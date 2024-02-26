# Client-side PKCE integration

This is a small Go program which serves up customised HTML and JavaScript pages,
to let clients exchange an OIDC PKCE code for a JWT token, which is then stored
in a Cookie.

It's been written with Kubernetes, Zitadel and Traefik in mind, but should work with any PKCE based Issuer.

## What's PKCE?

PKCE is an improvement to the (admittedly complex) mechanism by which a user logs in.

Before PKCE, it was possible for man-in-the-middle attacks to steal credentials
during the login process. PKCE makes this impossible, because it requires the client
to identify itself during the process of exchanging a login token for a bearer token
(which in my case is a JWT, but doesn't need to be).

Rather than do a poor job of explaining PKCE, here's some information written by
people who are probably smarter than I am:

* https://oauth.net/2/pkce/
* https://www.rfc-editor.org/rfc/rfc7636
* https://en.wikipedia.org/wiki/OAuth#Security_issues

## What does this do?

This code is a simple web server that serves up the JavaScript code requries to perform a generic PKCE exchange,
ultimately resuling in a [JWT bearer token](https://jwt.io/introduction/) being stored as a cookie in the browser.

When combined with an IDP like Zitadel, and a bit of middleware, this code lets you provide single-sign-on (SSO)
to all of your services, and to protect them behind a reverse proxy. In other words, using Zitadel, PKCE, and a proxy, 
users can log into your app using their Google, Apple or Microsoft accounts.

## Middleware

To make a service secure, we need a mechanism that checks that a HTTP request is allowed to access it.
In other words, when a request is made to your service's URL, we need to make sure the caller has permission to connect.
A common way to perform this check is by using JWT bearer tokens - which are just encoded cookies, sent
as part of the HTTP request. You can [read more about JWT bearer tokens here](https://jwt.io/introduction/).

You can certainly perform this check directly within your service endpoints, but in my option it's arguably more
secure and definitely easier to use JWT middleware such as the
[Brainnwave](https://github.com/Brainnwave/jwt-middleware) JWT middleware for Traefik. Brainnwave checks the user's
cookies for a JWT token, and redirects the browser to perform PKCE (via this module) if the JWT is missing or
invalid. A JWT might be invalid because it's expired, or perhaps because it has been forged.

To deploy Brainnwave for Traefik, a minimal K8s configuration for this combination might be:

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: jwt-brainnwave
spec:
    plugin:
      jwt:
        issuers:
          - https://issuer.example.com
        redirectUnauthorized: "https://hello.example.com/pkce/authorize.html?state={{.URL}}"
        redirectForbidden: "https://hello.example.com/unauthorized"
```

(you would also need to configure Traefik to use Brainnwave; I hope to post something about that later).

Brainnwave looks for the "Authorization" cookie by default. If this isn't detected, or if it's invalid, the user
is diverted to the `redirectUnauthorized` URL - which should be served by this application. That's all you need!

> Note that Brainnwave lets you configure all sorts of requirements for the JWT; you should
[read their docs](https://github.com/Brainnwave/jwt-middleware) to learn more.

## The role of PKCE

This module communicates with the Issuer to arrange for a JWT token to be stored in a cookie on the browser. That's it!
In general, it doesn't provide a user interface of any kind. It just loads up some JavaScript in the browser, performs
a few jumps through hoops, and sets the cookie.

Part of the "jumping through hoops" is redirecting the user to the IDP's login page. In other words, this code does
*not* log the user in. What it does is set up the conditions for a user to securely log in, and if successful,
downloads the JWT and puts it in a cookie - before sending the browser back to the user's application.