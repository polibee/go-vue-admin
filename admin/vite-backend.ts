export interface BackendEnv {
  VITE_BACKEND_PORT?: string
  VITE_BACKEND_URL?: string
}

export const defaultBackendPorts = ['3003', '3000', '3001', '3002']

type BackendProbe = (url: string) => Promise<boolean>

function normalizeUrl(url: string): string {
  return url.replace(/\/$/, '')
}

function candidateUrls(env: BackendEnv): string[] {
  const preferredPort = env.VITE_BACKEND_PORT?.trim()
  const ports = [...new Set([preferredPort, ...defaultBackendPorts].filter(Boolean))]
  return ports.map((port) => `http://127.0.0.1:${port}`)
}

async function probeBackend(url: string): Promise<boolean> {
  const controller = new AbortController()
  const timeout = setTimeout(() => controller.abort(), 400)
  try {
    const response = await fetch(`${url}/api/health`, {
      headers: { Accept: 'application/json' },
      signal: controller.signal,
    })
    return response.ok
  } catch {
    return false
  } finally {
    clearTimeout(timeout)
  }
}

export async function resolveBackendUrl(
  env: BackendEnv,
  probe: BackendProbe = probeBackend,
): Promise<string> {
  const configured = env.VITE_BACKEND_URL?.trim()
  if (configured) return normalizeUrl(configured)

  for (const candidate of candidateUrls(env)) {
    if (await probe(candidate)) return candidate
  }

  return 'http://127.0.0.1:3003'
}
