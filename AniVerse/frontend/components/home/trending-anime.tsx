'use client'

import { useQuery } from '@tanstack/react-query'
import { animeApi } from '@/lib/api'
import Link from 'next/link'
import { Star } from 'lucide-react'

export function TrendingAnime() {
  const { data: animes, isLoading } = useQuery({
    queryKey: ['trending'],
    queryFn: () => animeApi.getAll({ limit: 6 }).then(res => res.data),
  })

  if (isLoading) return <div className="py-12 text-center">Yukleniyor...</div>

  return (
    <section className="py-12 max-w-7xl mx-auto px-4">
      <h2 className="text-2xl font-bold mb-6">Trend Animeler</h2>
      
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
        {animes?.map((anime: any) => (
          <Link
            key={anime.id}
            href={`/anime/${anime.id}`}
            className="group"
          >
            <div className="relative aspect-[2/3] rounded-lg overflow-hidden mb-2">
              <img
                src={anime.cover_image || '/default-anime.png'}
                alt={anime.title}
                className="w-full h-full object-cover group-hover:scale-105 transition-transform"
              />
              <div className="absolute top-2 right-2 flex items-center gap-1 bg-black/60 px-2 py-1 rounded text-sm">
                <Star className="h-4 w-4 text-yellow-400 fill-yellow-400" />
                {anime.imdb_rating || 'N/A'}
              </div>
            </div>
            <h3 className="font-medium truncate group-hover:text-primary-400 transition-colors">
              {anime.title}
            </h3>
            <p className="text-sm text-gray-400">{anime.year}</p>
          </Link>
        ))}
      </div>
    </section>
  )
}