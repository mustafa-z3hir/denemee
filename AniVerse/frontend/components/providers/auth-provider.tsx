'use client'

import { useEffect } from 'react'
import { useAuthStore } from '@/lib/store'
import { userApi } from '@/lib/api'

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { setAuth, logout } = useAuthStore()

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      userApi
        .getMe()
        .then((res) => {
          setAuth(res.data, token)
        })
        .catch(() => {
          logout()
        })
    }
  }, [setAuth, logout])

  return <>{children}</>
}