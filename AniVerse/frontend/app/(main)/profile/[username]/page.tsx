'use client'

import { useParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import { userApi, badgeApi } from '@/lib/api'
import { Navbar } from '@/components/layout/navbar'
import { BadgeShowcase } from '@/components/badges/badge-showcase'

export default function ProfilePage() {
  const { username } = useParams()
  
  const { data: profile, isLoading: profileLoading } = useQuery({
    queryKey: ['profile', username],
    queryFn: () => userApi.getProfile(username as string).then(res => res.data),
  })

  const { data: badges, isLoading: badgesLoading } = useQuery({
    queryKey: ['badges', username],
    queryFn: () => badgeApi.getMyBadges().then(res => res.data),
    enabled: !!profile,
  })

  if (profileLoading) return <div className="min-h-screen bg-anime-dark p-8">Yukleniyor...</div>

  return (
    <div className="min-h-screen bg-anime-dark">
      <Navbar />
      
      <div className="h-64 bg-gradient-to-r from-primary-900 to-primary-700" />

      <div className="max-w-7xl mx-auto px-4 -mt-20">
        <div className="flex flex-col md:flex-row gap-8">
          <div className="flex-shrink-0">
            <img
              src={profile?.avatar || '/default-avatar.png'}
              alt={profile?.username}
              className="w-40 h-40 rounded-full border-4 border-anime-dark bg-anime-surface"
            />
          </div>

          <div className="flex-1 pt-4">
            <div className="flex items-center gap-3 mb-2">
              <h1 className="text-3xl font-bold">{profile?.username}</h1>
              {profile?.isVerified && (
                <span className="text-blue-400" title="Dogrulanmis">✓</span>
              )}
            </div>
            
            <p className="text-gray-400 mb-4">{profile?.bio || 'Biyografi yok'}</p>

            <div className="flex items-center gap-6 text-sm mb-6">
              <div>
                <span className="font-bold text-lg">{profile?.totalWatchHours || 0}</span>
                <span className="text-gray-400 ml-1">saat izleme</span>
              </div>
              <div>
                <span className="font-bold text-lg">{profile?.followersCount || 0}</span>
                <span className="text-gray-400 ml-1">takipci</span>
              </div>
              <div>
                <span className="font-bold text-lg">{profile?.followingCount || 0}</span>
                <span className="text-gray-400 ml-1">takip</span>
              </div>
            </div>
          </div>
        </div>

        <div className="mt-12">
          <h2 className="text-2xl font-bold mb-6">Rozetler</h2>
          {badgesLoading ? (
            <div>Yukleniyor...</div>
          ) : (
            <BadgeShowcase badges={badges || []} />
          )}
        </div>
      </div>
    </div>
  )
}