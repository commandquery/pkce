# Client-side PKCE integration

This is a small Go program which serves up HTML and JavaScript pages to let clients
exchance an OIDC PKCE code for a JWT token, which is then stored in a Cookie. It's been
written with Zitadel in mind, but should work with any PKCE based Issuer.

PKCE removes the need to store or provide a client secret to a browser in order for a user
to log in, since the client secret could be intercepted during login and used to obtain
credentials. Instead, the client generates a random verifier code, part of which
is sent to the server, and part of which is retained by the client and not initially sent.
Without the verifier code, a man-in-the-middle can't obtain the JWT token and, therefore
can't obtain the final JWT.

Here is more information about PKCE if you want it:

* https://oauth.net/2/pkce/
* https://www.rfc-editor.org/rfc/rfc7636
* https://en.wikipedia.org/wiki/OAuth#Security_issues

