import type { RequestEvent } from "@sveltejs/kit"

export const getUserToken = (event: RequestEvent) => {
	// get the cookies from the request
	const { cookies } = event

	// get the user token from the cookie
	const token = cookies.get("auth")


	if (token != undefined && token?.length > 0) {

        const user = {
			token: token
		}
		
		return user
	}
	
	return null
}