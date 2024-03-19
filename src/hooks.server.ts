import { getUserToken } from "$lib/server/auth"
import { redirect, type Handle } from "@sveltejs/kit"

export const handle: Handle = async ({ event, resolve }) => {
	
	event.locals.Token = getUserToken(event)

	if (event.url.pathname.startsWith("/auth")) {
		if (!event.locals.Token){
			throw redirect(303, "/")
		}
	
	}

	const response = await resolve(event) 


	return response
}