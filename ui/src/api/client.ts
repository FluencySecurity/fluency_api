/**
 * API client for the Fluency backend.
 * In dev, Vite proxies /api and /health to the backend (see vite.config.ts).
 */

const getBase = () => {
  // Use relative URLs so Vite proxy works in dev; in production you'd set VITE_API_BASE or serve from same origin
  if (import.meta.env.VITE_API_BASE) {
    return import.meta.env.VITE_API_BASE.replace(/\/$/, '')
  }
  return ''
}

export async function getHealth(signal?: AbortSignal): Promise<{ status: string }> {
  const res = await fetch(`${getBase()}/health`, { signal })
  if (!res.ok) throw new Error(`Health check failed: ${res.status}`)
  return res.json()
}

export interface EndpointEntry {
  method: string
  path: string
  description: string
}

export interface EndpointsResponse {
  endpoints: EndpointEntry[]
}

export async function getEndpoints(signal?: AbortSignal): Promise<EndpointsResponse> {
  const res = await fetch(`${getBase()}/api/endpoints`, { signal })
  if (!res.ok) throw new Error(`Failed to fetch endpoints: ${res.status}`)
  return res.json()
}

export interface SendRequestResult {
  status: number
  statusText: string
  body: unknown
  ok: boolean
}

/** Build URL with optional query params. */
function buildUrl(path: string, params?: Record<string, string>): string {
  const base = getBase()
  const url = `${base}${path.startsWith('/') ? path : `/${path}`}`
  if (!params || Object.keys(params).length === 0) return url
  const search = new URLSearchParams()
  for (const [k, v] of Object.entries(params)) {
    if (v !== '') search.set(k, v)
  }
  const q = search.toString()
  return q ? `${url}?${q}` : url
}

/**
 * Send an arbitrary request to the API. Body should be a string (e.g. JSON string) or undefined.
 */
export async function sendRequest(
  method: string,
  path: string,
  options?: { query?: Record<string, string>; body?: string }
): Promise<SendRequestResult> {
  const url = buildUrl(path, options?.query)
  const init: RequestInit = {
    method,
    headers: {},
  }
  if (options?.body !== undefined && options.body !== '') {
    init.headers = { 'Content-Type': 'application/json' }
    init.body = options.body
  }
  const res = await fetch(url, init)
  let body: unknown
  const ct = res.headers.get('content-type') || ''
  try {
    body = ct.includes('application/json') ? await res.json() : await res.text()
  } catch {
    body = null
  }
  return { status: res.status, statusText: res.statusText, body, ok: res.ok }
}
