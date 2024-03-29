import { SECRET_BASE_API_URL } from '$env/static/private'
import type { PageServerLoad } from "./$types"

interface Payment {
    id: string;
    tenant_id: string;
    amount: number;
    start_date: string;
    end_date: string;
    version: string;
    created_at: string;
    created_by: string;
    updated_at: string;
}

export const load : PageServerLoad = async ({ fetch, locals}): Promise<{ payments: Payment[] }>  => {

   let token = locals.Token

    const response = await fetch(`${SECRET_BASE_API_URL}/v1/auth/payments`, {
        method: "GET",
        headers: {
          "authorization": "Bearer " + token,
        },
    },
    )
    if (!response.ok) {
      return { payments : [] }
       
    }
    const payments: Payment[] = await response.json()
    return { payments : payments }
}