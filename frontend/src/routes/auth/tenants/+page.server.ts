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

export const load : PageServerLoad = async ({ fetch, locals}): Promise<{ tenants: Tenant[] }>  => {

   let token = locals.Token

    const response = await fetch(`${SECRET_BASE_API_URL}/v1/auth/tenants`, {
        method: "GET",
        headers: {
          "authorization": "Bearer " + token,
        },
    },
    )
    if (!response.ok) {
      return { tenants : [] }
       
    }
    const tenants: Tenant[] = await response.json()
    return { tenants : tenants }
}