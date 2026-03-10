'use client'

import { useEffect, useRef } from 'react'
import videojs from 'video.js'
import 'video.js/dist/video-js.css'

interface VideoPlayerProps {
  src: string
  poster?: string
  onProgress?: (time: number) => void
  onComplete?: () => void
}

export function VideoPlayer({ src, poster, onProgress, onComplete }: VideoPlayerProps) {
  const videoRef = useRef<HTMLDivElement>(null)
  const playerRef = useRef<any>(null)

  useEffect(() => {
    if (!videoRef.current) return

    playerRef.current = videojs(videoRef.current, {
      html5: {
        vhs: {
          overrideNative: true,
          limitRenditionByPlayerDimensions: true,
          useDevicePixelRatio: true,
        },
      },
      controls: true,
      fluid: true,
      poster,
      sources: [{ src, type: 'application/x-mpegURL' }],
    })

    playerRef.current.on('timeupdate', () => {
      onProgress?.(playerRef.current.currentTime())
    })

    playerRef.current.on('ended', () => {
      onComplete?.()
    })

    return () => {
      if (playerRef.current) {
        playerRef.current.dispose()
      }
    }
  }, [src])

  return (
    <div data-vjs-player>
      <video
        ref={videoRef}
        className="video-js vjs-big-play-centered vjs-theme-aniverse"
      />
    </div>
  )
}