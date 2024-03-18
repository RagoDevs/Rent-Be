import { SECRET_BASE_API_URL } from '$env/static/private'

export const load = async ({ fetch }) => {
  try {
    const response = await fetch(`${SECRET_BASE_API_URL}/v1/ping`)
    if (!response.ok) {
      throw new Error(`HTTP error: ${response.status}`)
    }
    const currencies = await response.json()
    return { currencies }
  } catch (error) {
    console.error(error)
    return { error: 'Unable to fetch currencies' }
  }
}