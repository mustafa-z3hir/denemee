import { Navbar } from '@/components/layout/navbar'
import { Footer } from '@/components/layout/footer'
import { Hero } from '@/components/home/hero'
import { TrendingAnime } from '@/components/home/trending-anime'

export default function Home() {
  return (
    <div className="min-h-screen bg-anime-dark">
      <Navbar />
      <main>
        <Hero />
        <TrendingAnime />
      </main>
      <Footer />
    </div>
  )
}