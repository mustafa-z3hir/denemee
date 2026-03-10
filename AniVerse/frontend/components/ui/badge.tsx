import { cn } from '@/lib/utils'

interface BadgeProps {
  children: React.ReactNode
  variant?: 'default' | 'secondary' | 'outline'
  className?: string
}

export function Badge({ children, variant = 'default', className }: BadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
        {
          'bg-primary-600 text-white': variant === 'default',
          'bg-anime-elevated text-gray-300': variant === 'secondary',
          'border border-anime-border text-gray-300': variant === 'outline',
        },
        className
      )}
    >
      {children}
    </span>
  )
}