//
// This is a POC wrapper library for the Fetch API that adds an access
// token to the header, and checks the result. If the result is a 403, we
// redirect to the login page to get the user to login again.
//

//
// TODO: life is better with TypeScript
//
async function call(input, init) {
    const accessToken = sessionStorage.getItem("bookworkAccessToken")
    if (accessToken == null) {
        window.location = "{{ .Base }}/login.html"
        return
    }

    const bearerToken = "Bearer " + accessToken

    // Add an Authorization token to any provided init parameters.
    if (init == null)
        init = { headers: { Authorization: bearerToken } }
    else if (init.headers == null)
        init.headers = { Authorization: bearerToken }
    else
        init.headers.Authorization = bearerToken

    const r = await fetch(input, init)
    if (r.status === 401)
        window.location = "{{ .Base }}/login.html"
    return r
}
