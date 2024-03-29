import { SECRET_BASE_API_URL } from '$env/static/private'
import type { PageServerLoad } from "./$types"

interface Payment {
    id: string;
    tenant_name: string;
    tenant_id: string;
    amount: number;
    start_date: string;
    end_date: string;
    admin_phone: string;
    location: string;
    block: string;
    partition: number;
    created_at: string;
    updated_at: string;
    version: string;
}



export const load : PageServerLoad = async ({ fetch, locals, params}): Promise<{ payment: Payment }>  => {

   let token = locals.Token

  let  uuid :string = params.uuid
   

    const response = await fetch(`${SECRET_BASE_API_URL}/v1/auth/payments/${uuid}`, {
        method: "GET",
        headers: {
          "authorization": "Bearer " + token,
        },
    },
    )
    if (!response.ok) {
      return { payment: {
                           id: '',
                           tenant_name: '',
                           tenant_id: '',
                           amount: 0,
                           start_date: '',
                           end_date: '',
                           admin_phone: '',
                           location: '',
                           block: '',
                           partition: 0,
                           created_at: '',
                           updated_at: '',
                           version: ''
          } };
       
    }
    const payment: Payment = await response.json()
    return { payment : payment }
}