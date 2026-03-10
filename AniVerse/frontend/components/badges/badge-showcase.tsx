'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'

interface Badge {
  id: string
  name: string
  icon: string
  color: string
  animation: string
  tier: number
  isEquipped: boolean
  description: string
  earnedAt: string
}

interface BadgeShowcaseProps {
  badges: Badge[]
  isEditable?: boolean
}

export function BadgeShowcase({ badges, isEditable = false }: BadgeShowcaseProps) {
  const [selectedBadge, setSelectedBadge] = useState<Badge | null>(null)
  const equipped = badges.filter(b => b.isEquipped).slice(0, 5)

  const getAnimationClass = (animation: string) => {
    switch (animation) {
      case 'pulse': return 'animate-pulse-slow'
      case 'glow': return 'animate-glow'
      case 'shimmer': return 'badge-shimmer'
      case 'rainbow': return 'animate-rainbow'
      default: return ''
    }
  }

  return (
    <div className="space-y-8">
      <div className="bg-anime-surface rounded-xl p-6 border border-anime-border">
        <h3 className="text-lg font-medium mb-4 text-gray-400">Ekipmanli Rozetler</h3>
        <div className="flex gap-4 flex-wrap">
          {equipped.map((badge, index) => (
            <motion.div
              key={badge.id}
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ delay: index * 0.1 }}
              onClick={() => setSelectedBadge(badge)}
              className={`
                relative w-20 h-20 rounded-xl flex items-center justify-center cursor-pointer
                transition-transform hover:scale-110
                ${getAnimationClass(badge.animation)}
              `}
              style={{ 
                backgroundColor: badge.color + '20',
                borderColor: badge.color,
                borderWidth: 2,
                boxShadow: `0 0 20px ${badge.color}40`
              }}
            >
              <span className="text-4xl">{badge.icon}</span>
              {badge.tier >= 8 && (
                <div className="absolute -top-1 -right-1 w-4 h-4 bg-yellow-400 rounded-full animate-pulse" />
              )}
            </motion.div>
          ))}
          
          {equipped.length < 3 && (
            <div className="w-20 h-20 rounded-xl border-2 border-dashed border-anime-border flex items-center justify-center text-gray-500">
              <span className="text-2xl">+</span>
            </div>
          )}
        </div>
      </div>

      <div>
        <h3 className="text-lg font-medium mb-4 text-gray-400">Tum Rozetler ({badges.length})</h3>
        <div className="grid grid-cols-4 sm:grid-cols-6 md:grid-cols-8 lg:grid-cols-10 gap-4">
          {badges.map((badge) => (
            <motion.div
              key={badge.id}
              whileHover={{ scale: 1.1 }}
              onClick={() => setSelectedBadge(badge)}
              className={`
                aspect-square rounded-lg flex items-center justify-center cursor-pointer
                ${badge.isEquipped ? 'ring-2 ring-primary-500' : 'opacity-60 hover:opacity-100'}
              `}
              style={{ backgroundColor: badge.color + '20' }}
            >
              <span className="text-2xl">{badge.icon}</span>
            </motion.div>
          ))}
        </div>
      </div>

      {selectedBadge && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-anime-surface rounded-2xl p-6 max-w-md w-full border border-anime-border"
          >
            <div className="flex items-center gap-4 mb-4">
              <div
                className={`w-20 h-20 rounded-xl flex items-center justify-center text-4xl ${getAnimationClass(selectedBadge.animation)}`}
                style={{ 
                  backgroundColor: selectedBadge.color + '30',
                  border: `2px solid ${selectedBadge.color}`
                }}
              >
                {selectedBadge.icon}
              </div>
              <div>
                <h3 className="text-xl font-bold">{selectedBadge.name}</h3>
                <p className="text-sm text-gray-400">Tier {selectedBadge.tier}</p>
              </div>
            </div>
            
            <p className="text-gray-300 mb-4">{selectedBadge.description}</p>
            
            <p className="text-sm text-gray-500 mb-6">
              Kazanildi: {new Date(selectedBadge.earnedAt).toLocaleDateString('tr-TR')}
            </p>

            <div className="flex gap-3">
              {isEditable && (
                <button
                  className="flex-1 py-2 bg-primary-600 hover:bg-primary-700 rounded-lg"
                  onClick={() => {}}
                >
                  {selectedBadge.isEquipped ? 'Cikar' : 'Ekipmanla'}
                </button>
              )}
              <button
                className="flex-1 py-2 bg-anime-elevated hover:bg-anime-border rounded-lg"
                onClick={() => setSelectedBadge(null)}
              >
                Kapat
              </button>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  )
}