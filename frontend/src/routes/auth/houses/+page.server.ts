import { SECRET_BASE_API_URL } from '$env/static/private'
import type { PageServerLoad } from "./$types"

interface House {
    id: string;
    location: string;
    block: string;
    partition: number;
    occupied: boolean;
}


export const load : PageServerLoad = async ({ fetch, locals}): Promise<{ houses: House[] }>  => {

   let token = locals.Token

    const response = await fetch(`${SECRET_BASE_API_URL}/v1/auth/houses`, {
        method: "GET",
        headers: {
          "authorization": "Bearer " + token,
        },
    },
    )
    if (!response.ok) {
      return { houses : [] }
       
    }
    const houses: House[] = await response.json()
    return { houses : houses }
}