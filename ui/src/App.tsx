import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Dashboard from './pages/Dashboard'
import RequestBuilder from './pages/RequestBuilder'
import './App.css'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/request" element={<RequestBuilder />} />
      </Routes>
    </BrowserRouter>
  )
}
