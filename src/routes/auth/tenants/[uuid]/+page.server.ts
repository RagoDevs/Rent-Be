import { SECRET_BASE_API_URL } from '$env/static/private'
import type { PageServerLoad } from "./$types"

interface Tenant {
    id: string;
    first_name: string;
    last_name: string;
    house_id: string;
    phone: string;
    personal_id_type: string;
    personal_id: string;
    active: boolean;
    sos: string; 
}

export const load : PageServerLoad = async ({ fetch, locals, params}): Promise<{ tenant: Tenant }>  => {

   let token = locals.Token
   let  uuid :string = params.uuid
   

    const response = await fetch(`${SECRET_BASE_API_URL}/v1/auth/tenants/${uuid}`, {
        method: "GET",
        headers: {
          "authorization": "Bearer " + token,
        },
    },
    )
    if (!response.ok) {
      return { tenant : {id:"", first_name: "", last_name : "", house_id : "",phone : "",personal_id: "", personal_id_type: "", active : false ,sos : ""}, }
       
    }
    const tenant: Tenant = await response.json()
    return { tenant : tenant }
}