import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { getHealth, getEndpoints, type EndpointEntry } from '../api/client'
import '../App.css'

export default function Dashboard() {
  const [health, setHealth] = useState<{ status: string } | null>(null)
  const [endpoints, setEndpoints] = useState<EndpointEntry[]>([])
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const controller = new AbortController()
    const signal = controller.signal
    let cancelled = false
    setLoading(true)
    setError(null)

    Promise.all([getHealth(signal), getEndpoints(signal)])
      .then(([h, { endpoints: list }]) => {
        if (!cancelled) {
          setHealth(h)
          setEndpoints(list)
        }
      })
      .catch((e) => {
        if (!cancelled && e?.name !== 'AbortError') setError(e instanceof Error ? e.message : String(e))
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })

    return () => {
      cancelled = true
      controller.abort()
    }
  }, [])

  if (loading) {
    return (
      <div className="app">
        <h1>Fluency</h1>
        <p className="muted">Connecting to API…</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="app">
        <h1>Fluency</h1>
        <p className="error">{error}</p>
        <p className="muted">Ensure the backend is running: <code>./fluency serve</code></p>
      </div>
    )
  }

  return (
    <div className="app">
      <header>
        <h1>Fluency</h1>
        <p className="muted">API dashboard</p>
      </header>

      <section className="card">
        <h2>Health</h2>
        <p className="status">{health?.status === 'ok' ? '✓ OK' : health?.status ?? '—'}</p>
      </section>

      <section className="card">
        <h2>Endpoints</h2>
        <p className="muted">{endpoints.length} routes</p>
        <ul className="endpoint-list">
          {endpoints.map((ep) => (
            <li key={`${ep.method} ${ep.path}`}>
              <Link
                to="/request"
                state={{ endpoint: ep }}
                className="endpoint-link"
              >
                <span className="method">{ep.method}</span>
                <code>{ep.path}</code>
                <span className="desc">{ep.description}</span>
              </Link>
            </li>
          ))}
        </ul>
      </section>
    </div>
  )
}
