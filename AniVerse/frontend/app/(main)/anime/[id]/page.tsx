'use client'

import { useParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import { animeApi } from '@/lib/api'
import { Navbar } from '@/components/layout/navbar'
import { Badge } from '@/components/ui/badge'
import { Star, Clock, Calendar, Play } from 'lucide-react'

export default function AnimePage() {
  const { id } = useParams()
  
  const { data: anime, isLoading } = useQuery({
    queryKey: ['anime', id],
    queryFn: () => animeApi.getById(id as string).then(res => res.data),
  })

  if (isLoading) return <div className="min-h-screen bg-anime-dark p-8">Yukleniyor...</div>
  if (!anime) return <div className="min-h-screen bg-anime-dark p-8">Anime bulunamadi</div>

  return (
    <div className="min-h-screen bg-anime-dark">
      <Navbar />
      
      <div className="relative h-96">
        <img
          src={anime.banner_image || anime.cover_image}
          alt={anime.title}
          className="w-full h-full object-cover"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-anime-dark via-anime-dark/50 to-transparent" />
      </div>

      <div className="max-w-7xl mx-auto px-4 -mt-32 relative z-10">
        <div className="flex flex-col md:flex-row gap-8">
          <img
            src={anime.cover_image}
            alt={anime.title}
            className="w-64 rounded-lg shadow-2xl"
          />

          <div className="flex-1">
            <h1 className="text-4xl font-bold mb-2">{anime.title}</h1>
            <h2 className="text-xl text-gray-400 mb-4">{anime.title_en}</h2>

            <div className="flex flex-wrap gap-2 mb-4">
              {anime.genres?.map((genre: any) => (
                <Badge key={genre.id} variant="secondary">
                  {genre.name}
                </Badge>
              ))}
            </div>

            <div className="flex items-center gap-6 mb-6">
              <div className="flex items-center gap-1">
                <Star className="h-5 w-5 text-yellow-400 fill-yellow-400" />
                <span className="text-lg font-bold">{anime.imdb_rating || 'N/A'}</span>
                <span className="text-gray-400">IMDb</span>
              </div>
              <div className="flex items-center gap-1">
                <Clock className="h-5 w-5 text-gray-400" />
                <span>{anime.duration} dk/bolum</span>
              </div>
              <div className="flex items-center gap-1">
                <Calendar className="h-5 w-5 text-gray-400" />
                <span>{anime.year}</span>
              </div>
            </div>

            <p className="text-gray-300 mb-6 line-clamp-4">{anime.description}</p>

            <button className="flex items-center gap-2 px-6 py-3 bg-primary-600 hover:bg-primary-700 rounded-lg font-medium">
              <Play className="h-5 w-5" />
              Izlemeye Basla
            </button>
          </div>
        </div>

        <div className="mt-12">
          <h3 className="text-2xl font-bold mb-4">Bolumler</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
            {anime.episodes_list?.map((episode: any) => (
              <div
                key={episode.id}
                className="bg-anime-surface rounded-lg overflow-hidden hover:ring-2 ring-primary-500 cursor-pointer transition-all"
              >
                <img
                  src={episode.thumbnail || anime.cover_image}
                  alt={`Bolum ${episode.number}`}
                  className="w-full aspect-video object-cover"
                />
                <div className="p-3">
                  <p className="font-medium">Bolum {episode.number}</p>
                  <p className="text-sm text-gray-400 truncate">{episode.title}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}