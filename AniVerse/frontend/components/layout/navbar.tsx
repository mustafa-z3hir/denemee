'use client'

import Link from 'next/link'
import { useAuthStore } from '@/lib/store'
import { Search, Bell, MessageSquare, User, Menu, X } from 'lucide-react'
import { useState } from 'react'

export function Navbar() {
  const { user, isAuthenticated, logout } = useAuthStore()
  const [isMenuOpen, setIsMenuOpen] = useState(false)

  return (
    <nav className="sticky top-0 z-50 bg-anime-surface/80 backdrop-blur-md border-b border-anime-border">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <Link href="/" className="flex items-center space-x-2">
            <span className="text-2xl font-bold bg-gradient-to-r from-primary-400 to-primary-600 bg-clip-text text-transparent">
              AniVerse
            </span>
          </Link>

          <div className="hidden md:flex flex-1 max-w-md mx-8">
            <div className="relative w-full">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
              <input
                type="text"
                placeholder="Anime ara..."
                className="w-full pl-10 pr-4 py-2 bg-anime-elevated border border-anime-border rounded-lg text-sm focus:outline-none focus:border-primary-500"
              />
            </div>
          </div>

          <div className="flex items-center space-x-4">
            {isAuthenticated ? (
              <>
                <button className="p-2 hover:bg-anime-elevated rounded-lg relative">
                  <Bell className="h-5 w-5 text-gray-400" />
                  <span className="absolute top-1 right-1 h-2 w-2 bg-red-500 rounded-full" />
                </button>
                
                <Link href="/messages" className="p-2 hover:bg-anime-elevated rounded-lg">
                  <MessageSquare className="h-5 w-5 text-gray-400" />
                </Link>

                <div className="relative group">
                  <button className="flex items-center space-x-2 p-2 hover:bg-anime-elevated rounded-lg">
                    <img
                      src={user?.avatar || '/default-avatar.png'}
                      alt={user?.username}
                      className="h-8 w-8 rounded-full"
                    />
                    <span className="hidden sm:block text-sm font-medium">
                      {user?.username}
                    </span>
                  </button>

                  <div className="absolute right-0 mt-2 w-48 bg-anime-surface border border-anime-border rounded-lg shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all">
                    <Link href={`/profile/${user?.username}`} className="block px-4 py-2 hover:bg-anime-elevated">
                      Profilim
                    </Link>
                    <Link href="/settings" className="block px-4 py-2 hover:bg-anime-elevated">
                      Ayarlar
                    </Link>
                    {user?.isAdmin && (
                      <Link href="/admin" className="block px-4 py-2 hover:bg-anime-elevated text-primary-400">
                        Admin Panel
                      </Link>
                    )}
                    <hr className="border-anime-border" />
                    <button
                      onClick={logout}
                      className="block w-full text-left px-4 py-2 hover:bg-anime-elevated text-red-400"
                    >
                      Cikis Yap
                    </button>
                  </div>
                </div>
              </>
            ) : (
              <div className="flex items-center space-x-2">
                <Link
                  href="/login"
                  className="px-4 py-2 text-sm font-medium hover:text-primary-400"
                >
                  Giris
                </Link>
                <Link
                  href="/register"
                  className="px-4 py-2 text-sm font-medium bg-primary-600 hover:bg-primary-700 rounded-lg"
                >
                  Kayit Ol
                </Link>
              </div>
            )}

            <button
              className="md:hidden p-2"
              onClick={() => setIsMenuOpen(!isMenuOpen)}
            >
              {isMenuOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
            </button>
          </div>
        </div>
      </div>

      {isMenuOpen && (
        <div className="md:hidden bg-anime-surface border-t border-anime-border">
          <div className="px-4 py-3">
            <input
              type="text"
              placeholder="Anime ara..."
              className="w-full px-4 py-2 bg-anime-elevated border border-anime-border rounded-lg"
            />
          </div>
        </div>
      )}
    </nav>
  )
}