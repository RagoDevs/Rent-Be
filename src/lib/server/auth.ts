import type { RequestEvent } from "@sveltejs/kit"

export const getUserToken = (event: RequestEvent) => {
	// get the cookies from the request
	const { cookies } = event

	let token: string | null;

    const authCookie = cookies.get("auth");
    if (authCookie !== undefined) {
        token = authCookie;
    } else {
        token = null;
    }


	if (token != undefined && token?.length > 0) {

        const Token = token
		return Token
	}
	
	return null
}