import { useState, useCallback } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { sendRequest } from '../api/client'
import type { EndpointEntry } from '../api/client'
import '../App.css'
import './RequestBuilder.css'

type ParamRow = { key: string; value: string }

function parseQueryHint(description: string): string[] {
  const m = description.match(/\(query:\s*([^)]+)\)/i)
  if (!m) return []
  return m[1].split(',').map((s) => s.trim())
}

export default function RequestBuilder() {
  const { state } = useLocation() as { state?: { endpoint: EndpointEntry } }
  const endpoint = state?.endpoint

  const [queryParams, setQueryParams] = useState<ParamRow[]>(() => {
    if (!endpoint) return []
    const hints = parseQueryHint(endpoint.description)
    return hints.length ? hints.map((k) => ({ key: k, value: '' })) : [{ key: '', value: '' }]
  })
  const [body, setBody] = useState('')
  const [sending, setSending] = useState(false)
  const [result, setResult] = useState<{
    status: number
    statusText: string
    body: unknown
    ok: boolean
  } | null>(null)
  const [error, setError] = useState<string | null>(null)

  const addQueryRow = useCallback(() => {
    setQueryParams((prev) => [...prev, { key: '', value: '' }])
  }, [])

  const removeQueryRow = useCallback((i: number) => {
    setQueryParams((prev) => prev.filter((_, idx) => idx !== i))
  }, [])

  const updateQueryRow = useCallback((i: number, field: 'key' | 'value', value: string) => {
    setQueryParams((prev) => {
      const next = [...prev]
      next[i] = { ...next[i], [field]: value }
      return next
    })
  }, [])

  const buildQuery = useCallback((): Record<string, string> => {
    const out: Record<string, string> = {}
    for (const row of queryParams) {
      if (row.key.trim()) out[row.key.trim()] = row.value
    }
    return out
  }, [queryParams])

  const handleSend = useCallback(async () => {
    if (!endpoint) return
    setSending(true)
    setError(null)
    setResult(null)
    try {
      const query = buildQuery()
      const bodyStr = body.trim() || undefined
      const res = await sendRequest(endpoint.method, endpoint.path, {
        query: Object.keys(query).length ? query : undefined,
        body: bodyStr,
      })
      setResult(res)
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
    } finally {
      setSending(false)
    }
  }, [endpoint, buildQuery, body])

  if (!endpoint) {
    return (
      <div className="app">
        <p className="muted">No endpoint selected.</p>
        <Link to="/">Back to dashboard</Link>
      </div>
    )
  }

  const hasBody = ['POST', 'PUT', 'PATCH'].includes(endpoint.method) ||
    endpoint.path.includes('search') || endpoint.path.includes('eventwatch')

  return (
    <div className="app">
      <header className="request-builder-header">
        <Link to="/" className="back-link">← Dashboard</Link>
        <h1 className="request-title">
          <span className="method-badge">{endpoint.method}</span>
          <code>{endpoint.path}</code>
        </h1>
        <p className="muted">{endpoint.description}</p>
      </header>

      <section className="card request-builder-card">
        <h2>Query parameters</h2>
        {queryParams.map((row, i) => (
          <div key={i} className="param-row">
            <input
              type="text"
              placeholder="name"
              value={row.key}
              onChange={(e) => updateQueryRow(i, 'key', e.target.value)}
              className="param-key"
            />
            <input
              type="text"
              placeholder="value"
              value={row.value}
              onChange={(e) => updateQueryRow(i, 'value', e.target.value)}
              className="param-value"
            />
            <button type="button" onClick={() => removeQueryRow(i)} className="param-remove">
              Remove
            </button>
          </div>
        ))}
        <button type="button" onClick={addQueryRow} className="param-add">+ Add query param</button>
      </section>

      {hasBody && (
        <section className="card request-builder-card">
          <h2>Body (JSON)</h2>
          <textarea
            value={body}
            onChange={(e) => setBody(e.target.value)}
            placeholder='{"query": "..."}'
            className="body-textarea"
            rows={8}
          />
        </section>
      )}

      <div className="send-row">
        <button onClick={handleSend} disabled={sending} className="send-btn">
          {sending ? 'Sending…' : 'Send request'}
        </button>
      </div>

      {error && (
        <section className="card response-card error-card">
          <h2>Error</h2>
          <pre>{error}</pre>
        </section>
      )}

      {result && (
        <section className={`card response-card ${result.ok ? '' : 'error-card'}`}>
          <h2>Response</h2>
          <p className="response-status">
            {result.status} {result.statusText}
          </p>
          <pre className="response-body">
            {typeof result.body === 'string'
              ? result.body
              : JSON.stringify(result.body, null, 2)}
          </pre>
        </section>
      )}
    </div>
  )
}
