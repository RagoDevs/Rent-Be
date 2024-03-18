import { getUserToken } from "$lib/server/auth"
import { redirect, type Handle } from "@sveltejs/kit"

export const handle: Handle = async ({ event, resolve }) => {
	
	event.locals.user = getUserToken(event)

	if (event.url.pathname.startsWith("/auth")) {
		if (!event.locals.user?.token){
			throw redirect(303, "/")
		}
	
	}

	const response = await resolve(event) 


	return response
}