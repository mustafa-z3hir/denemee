'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { authApi } from '@/lib/api'
import { useAuthStore } from '@/lib/store'

export default function LoginPage() {
  const router = useRouter()
  const { setAuth } = useAuthStore()
  const [form, setForm] = useState({ email: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    try {
      const res = await authApi.login(form)
      setAuth(res.data.user, res.data.token)
      router.push('/')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Giris basarisiz')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-anime-dark px-4">
      <div className="max-w-md w-full space-y-8 bg-anime-surface p-8 rounded-2xl border border-anime-border">
        <div className="text-center">
          <h2 className="text-3xl font-bold bg-gradient-to-r from-primary-400 to-primary-600 bg-clip-text text-transparent">
            AniVerse
          </h2>
          <p className="mt-2 text-gray-400">Hesabina giris yap</p>
        </div>

        {error && (
          <div className="bg-red-500/10 border border-red-500 text-red-500 px-4 py-2 rounded-lg text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-medium mb-2">Email</label>
            <input
              type="email"
              required
              className="w-full px-4 py-3 bg-anime-elevated border border-anime-border rounded-lg focus:outline-none focus:border-primary-500"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Sifre</label>
            <input
              type="password"
              required
              className="w-full px-4 py-3 bg-anime-elevated border border-anime-border rounded-lg focus:outline-none focus:border-primary-500"
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 bg-primary-600 hover:bg-primary-700 rounded-lg font-medium transition-colors disabled:opacity-50"
          >
            {loading ? 'Giris yapiliyor...' : 'Giris Yap'}
          </button>
        </form>

        <p className="text-center text-sm text-gray-400">
          Hesabin yok mu?{' '}
          <Link href="/register" className="text-primary-400 hover:underline">
            Kayit ol
          </Link>
        </p>
      </div>
    </div>
  )
}