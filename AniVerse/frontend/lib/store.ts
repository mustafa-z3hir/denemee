import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface User {
  id: string
  email: string
  username: string
  avatar?: string
  isAdmin: boolean
  isVerified: boolean
  equippedSlots: number
  maxSlots: number
  aniPoints: number
  totalWatchHours: number
}

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  setAuth: (user: User, token: string) => void
  logout: () => void
  updateUser: (user: Partial<User>) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      setAuth: (user, token) => {
        localStorage.setItem('token', token)
        set({ user, token, isAuthenticated: true })
      },
      logout: () => {
        localStorage.removeItem('token')
        set({ user: null, token: null, isAuthenticated: false })
      },
      updateUser: (userData) =>
        set((state) => ({
          user: state.user ? { ...state.user, ...userData } : null,
        })),
    }),
    {
      name: 'auth-storage',
    }
  )
)

interface Badge {
  id: string
  type: string
  category: string
  tier: number
  name: string
  icon: string
  color: string
  animation: string
  description: string
  isEquipped: boolean
  slot?: number
  earnedAt: string
  earnedHow: string
}

interface BadgeState {
  badges: Badge[]
  equipped: Badge[]
  setBadges: (badges: Badge[]) => void
  equipBadge: (badgeId: string, slot: number) => void
  unequipBadge: (badgeId: string) => void
}

export const useBadgeStore = create<BadgeState>((set) => ({
  badges: [],
  equipped: [],
  setBadges: (badges) =>
    set({
      badges,
      equipped: badges.filter((b) => b.isEquipped).sort((a, b) => (a.slot || 0) - (b.slot || 0)),
    }),
  equipBadge: (badgeId, slot) =>
    set((state) => {
      const newBadges = state.badges.map((b) =>
        b.id === badgeId ? { ...b, isEquipped: true, slot } : b
      )
      return {
        badges: newBadges,
        equipped: newBadges.filter((b) => b.isEquipped).sort((a, b) => (a.slot || 0) - (b.slot || 0)),
      }
    }),
  unequipBadge: (badgeId) =>
    set((state) => {
      const newBadges = state.badges.map((b) =>
        b.id === badgeId ? { ...b, isEquipped: false, slot: undefined } : b
      )
      return {
        badges: newBadges,
        equipped: newBadges.filter((b) => b.isEquipped).sort((a, b) => (a.slot || 0) - (b.slot || 0)),
      }
    }),
}))

interface SocketState {
  socket: any | null
  isConnected: boolean
  setSocket: (socket: any) => void
  setConnected: (connected: boolean) => void
}

export const useSocketStore = create<SocketState>((set) => ({
  socket: null,
  isConnected: false,
  setSocket: (socket) => set({ socket }),
  setConnected: (isConnected) => set({ isConnected }),
}))