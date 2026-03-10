'use client'

import Link from 'next/link'
import { Play, Info } from 'lucide-react'

export function Hero() {
  return (
    <div className="relative h-[70vh] flex items-center">
      <div className="absolute inset-0 bg-gradient-to-r from-primary-900/50 to-anime-dark" />
      
      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="max-w-2xl">
          <h1 className="text-5xl md:text-6xl font-bold mb-4">
            Anime Dunyasini Kesfet
          </h1>
          <p className="text-xl text-gray-300 mb-8">
            En populer animeleri izle, topluluga katil, rozetler topla ve arkadaslarinla paylas.
          </p>
          <div className="flex gap-4">
            <Link
              href="/anime"
              className="flex items-center gap-2 px-6 py-3 bg-primary-600 hover:bg-primary-700 rounded-lg font-medium"
            >
              <Play className="h-5 w-5" />
              Izlemeye Basla
            </Link>
            <Link
              href="/fansubs"
              className="flex items-center gap-2 px-6 py-3 bg-anime-elevated hover:bg-anime-border rounded-lg font-medium"
            >
              <Info className="h-5 w-5" />
              Fansublar
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}