'use client'

import { useEffect, useState, useRef } from 'react'
import { Navbar } from '@/components/layout/navbar'
import { useSocketStore } from '@/lib/store'
import { dmApi } from '@/lib/api'
import { useQuery } from '@tanstack/react-query'
import { Send, Paperclip, Image, Smile } from 'lucide-react'

export default function MessagesPage() {
  const { socket } = useSocketStore()
  const [conversations, setConversations] = useState<any[]>([])
  const [activeConversation, setActiveConversation] = useState<string | null>(null)
  const [messages, setMessages] = useState<any[]>([])
  const [newMessage, setNewMessage] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const { data: convs } = useQuery({
    queryKey: ['conversations'],
    queryFn: () => dmApi.getConversations().then(res => res.data),
  })

  useEffect(() => {
    if (convs) setConversations(convs)
  }, [convs])

  useEffect(() => {
    if (!socket) return

    socket.on('message', (msg: any) => {
      setMessages(prev => [...prev, msg])
    })

    return () => {
      socket.off('message')
    }
  }, [socket])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const loadMessages = async (conversationId: string) => {
    setActiveConversation(conversationId)
    const res = await dmApi.getMessages(conversationId)
    setMessages(res.data)
  }

  const sendMessage = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newMessage.trim() || !activeConversation) return

    await dmApi.sendMessage(activeConversation, { content: newMessage })
    setNewMessage('')
  }

  return (
    <div className="min-h-screen bg-anime-dark">
      <Navbar />
      
      <div className="flex h-[calc(100vh-64px)]">
        <div className="w-80 bg-anime-surface border-r border-anime-border">
          <div className="p-4 border-b border-anime-border">
            <h2 className="text-lg font-bold">Mesajlar</h2>
          </div>
          <div className="overflow-y-auto">
            {conversations.map((conv) => (
              <button
                key={conv.id}
                onClick={() => loadMessages(conv.id)}
                className={`w-full p-4 flex items-center gap-3 hover:bg-anime-elevated transition-colors ${
                  activeConversation === conv.id ? 'bg-anime-elevated' : ''
                }`}
              >
                <img
                  src={conv.participants?.[0]?.user?.avatar || '/default-avatar.png'}
                  alt=""
                  className="w-12 h-12 rounded-full"
                />
                <div className="flex-1 text-left">
                  <p className="font-medium truncate">
                    {conv.type === 'direct' 
                      ? conv.participants?.[0]?.user?.username 
                      : conv.title}
                  </p>
                  <p className="text-sm text-gray-400 truncate">
                    {conv.messages?.[0]?.content || 'Henuz mesaj yok'}
                  </p>
                </div>
              </button>
            ))}
          </div>
        </div>

        <div className="flex-1 flex flex-col">
          {activeConversation ? (
            <>
              <div className="flex-1 overflow-y-auto p-4 space-y-4">
                {messages.map((msg, idx) => (
                  <div
                    key={idx}
                    className={`flex ${msg.sender_id === 'me' ? 'justify-end' : 'justify-start'}`}
                  >
                    <div
                      className={`max-w-[70%] px-4 py-2 rounded-2xl ${
                        msg.sender_id === 'me'
                          ? 'bg-primary-600 rounded-br-none'
                          : 'bg-anime-elevated rounded-bl-none'
                      }`}
                    >
                      <p>{msg.content}</p>
                      <span className="text-xs opacity-60">
                        {new Date(msg.created_at).toLocaleTimeString('tr-TR', { 
                          hour: '2-digit', 
                          minute: '2-digit' 
                        })}
                      </span>
                    </div>
                  </div>
                ))}
                <div ref={messagesEndRef} />
              </div>

              <form onSubmit={sendMessage} className="p-4 bg-anime-surface border-t border-anime-border">
                <div className="flex items-center gap-2">
                  <button type="button" className="p-2 hover:bg-anime-elevated rounded-lg">
                    <Paperclip className="h-5 w-5 text-gray-400" />
                  </button>
                  <button type="button" className="p-2 hover:bg-anime-elevated rounded-lg">
                    <Image className="h-5 w-5 text-gray-400" />
                  </button>
                  <input
                    type="text"
                    value={newMessage}
                    onChange={(e) => setNewMessage(e.target.value)}
                    placeholder="Mesaj yaz..."
                    className="flex-1 px-4 py-2 bg-anime-elevated border border-anime-border rounded-lg focus:outline-none focus:border-primary-500"
                  />
                  <button type="button" className="p-2 hover:bg-anime-elevated rounded-lg">
                    <Smile className="h-5 w-5 text-gray-400" />
                  </button>
                  <button
                    type="submit"
                    className="p-2 bg-primary-600 hover:bg-primary-700 rounded-lg"
                  >
                    <Send className="h-5 w-5" />
                  </button>
                </div>
              </form>
            </>
          ) : (
            <div className="flex-1 flex items-center justify-center text-gray-400">
              Bir sohbet secin
            </div>
          )}
        </div>
      </div>
    </div>
  )
}