'use client'

import { useAuthStore } from '@/lib/store'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { adminApi } from '@/lib/api'
import { useQuery } from '@tanstack/react-query'

export default function AdminPage() {
  const { user } = useAuthStore()
  const router = useRouter()

  useEffect(() => {
    if (!user?.isAdmin) {
      router.push('/')
    }
  }, [user, router])

  const { data: stats } = useQuery({
    queryKey: ['admin-dashboard'],
    queryFn: () => adminApi.getDashboard().then(res => res.data),
    enabled: user?.isAdmin,
  })

  if (!user?.isAdmin) return null

  return (
    <div className="min-h-screen bg-anime-dark p-8">
      <h1 className="text-3xl font-bold mb-8">Admin Panel</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="bg-anime-surface p-6 rounded-xl border border-anime-border">
          <p className="text-gray-400 text-sm">Toplam Kullanici</p>
          <p className="text-3xl font-bold">{stats?.totalUsers || 0}</p>
        </div>
        <div className="bg-anime-surface p-6 rounded-xl border border-anime-border">
          <p className="text-gray-400 text-sm">Bugun Aktif</p>
          <p className="text-3xl font-bold">{stats?.activeToday || 0}</p>
        </div>
        <div className="bg-anime-surface p-6 rounded-xl border border-anime-border">
          <p className="text-gray-400 text-sm">Toplam Anime</p>
          <p className="text-3xl font-bold">{stats?.totalAnime || 0}</p>
        </div>
        <div className="bg-anime-surface p-6 rounded-xl border border-anime-border">
          <p className="text-gray-400 text-sm">Bekleyen Rapor</p>
          <p className="text-3xl font-bold text-yellow-400">{stats?.pendingReports || 0}</p>
        </div>
      </div>

      <div className="bg-anime-surface rounded-xl border border-anime-border overflow-hidden">
        <div className="p-4 border-b border-anime-border">
          <h2 className="text-lg font-bold">Son Kullanicilar</h2>
        </div>
        <table className="w-full">
          <thead className="bg-anime-elevated">
            <tr>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Kullanici</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Email</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Kayit Tarihi</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-400">Islemler</th>
            </tr>
          </thead>
          <tbody>
            {stats?.recentUsers?.map((u: any) => (
              <tr key={u.id} className="border-t border-anime-border">
                <td className="px-4 py-3">{u.username}</td>
                <td className="px-4 py-3 text-gray-400">{u.email}</td>
                <td className="px-4 py-3 text-gray-400">
                  {new Date(u.created_at).toLocaleDateString('tr-TR')}
                </td>
                <td className="px-4 py-3">
                  <button className="text-red-400 hover:text-red-300 text-sm">
                    Banla
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}