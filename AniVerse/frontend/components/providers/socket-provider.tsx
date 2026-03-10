'use client'

import { useEffect } from 'react'
import { io } from 'socket.io-client'
import { useAuthStore } from '@/lib/store'
import { useSocketStore } from '@/lib/store'

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080'

export function SocketProvider({ children }: { children: React.ReactNode }) {
  const { token } = useAuthStore()
  const { setSocket, setConnected } = useSocketStore()

  useEffect(() => {
    if (!token) return

    const socket = io(WS_URL, {
      auth: { token },
      transports: ['websocket'],
    })

    socket.on('connect', () => {
      console.log('Socket connected')
      setConnected(true)
    })

    socket.on('disconnect', () => {
      console.log('Socket disconnected')
      setConnected(false)
    })

    socket.on('message', (data: any) => {
      console.log('New message:', data)
    })

    setSocket(socket)

    return () => {
      socket.disconnect()
    }
  }, [token, setSocket, setConnected])

  return <>{children}</>
}