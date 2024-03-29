import { redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";
import { SECRET_BASE_API_URL } from '$env/static/private'

export const actions: Actions = {
  default: async ({ cookies, request }) => {

    const data = await request.formData();
		const phone = data.get('phone');
		const password = data.get('password');


      const credentials = {
        phone: phone,
        password: password,
      };

      // Sending a POST request to the login URL
      const response = await fetch(`${SECRET_BASE_API_URL}/v1/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(credentials),
      });

      if (!response.ok) {
        // throw new Error("Login failed");
        return { success: false , error : "login failed"};
      }

    
      const token = await response.json();

      // Set the token in the cookies
      cookies.set("auth", token.token, {
        path: "/",
        httpOnly: true,
        sameSite: "strict",
        secure: process.env.NODE_ENV === "production",
        maxAge: 60 * 60 * 24 , // 1 day
      });

      return redirect(303, "/auth/dashboard");
    
  },
};
