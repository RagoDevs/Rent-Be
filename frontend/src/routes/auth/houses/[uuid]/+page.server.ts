import { SECRET_BASE_API_URL } from '$env/static/private'
import type { PageServerLoad } from "./$types"

interface House {
    id: string;
    location: string;
    block: string;
    partition: number;
    occupied: boolean;
}


export const load : PageServerLoad = async ({ fetch, locals, params}): Promise<{ house: House }>  => {

   let token = locals.Token

  let  uuid :string = params.uuid
   

    const response = await fetch(`${SECRET_BASE_API_URL}/v1/auth/houses/${uuid}`, {
        method: "GET",
        headers: {
          "authorization": "Bearer " + token,
        },
    },
    )
    if (!response.ok) {
      return { house: { id: '', location: '', block: '', partition: 0, occupied: false } };
       
    }
    const house: House = await response.json()
    return { house : house }
}