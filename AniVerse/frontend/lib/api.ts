import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export const api = axios.create({
  baseURL: `${API_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const authApi = {
  register: (data: { email: string; username: string; password: string }) =>
    api.post('/auth/register', data),
  
  login: (data: { email: string; password: string }) =>
    api.post('/auth/login', data),
  
  refresh: () => api.post('/auth/refresh'),
}

export const animeApi = {
  getAll: (params?: { page?: number; limit?: number; title?: string; genre?: string; status?: string }) =>
    api.get('/anime', { params }),
  
  getById: (id: string) => api.get(`/anime/${id}`),
  
  getEpisodes: (id: string) => api.get(`/anime/${id}/episodes`),
  
  rate: (id: string, data: { score: number; review?: string; isSpoiler?: boolean }) =>
    api.post(`/anime/${id}/rate`, data),
  
  getRatings: (id: string) => api.get(`/anime/${id}/ratings`),
  
  recordWatch: (episodeId: string, data: { progress: number; isCompleted: boolean }) =>
    api.post(`/anime/${episodeId}/watch`, data),
}

export const userApi = {
  getMe: () => api.get('/user/me'),
  
  getProfile: (username: string) => api.get(`/user/${username}`),
  
  updateProfile: (data: Partial<{ bio: string; avatar: string }>) =>
    api.put('/user/me', data),
  
  follow: (id: string) => api.post(`/user/follow/${id}`),
  
  getHistory: () => api.get('/history'),
  
  getContinueWatching: () => api.get('/history/continue'),
}

export const badgeApi = {
  getMyBadges: () => api.get('/badges'),
  
  checkNew: () => api.get('/badges/check'),
  
  equip: (data: { badgeId: string; equipped: boolean; slot?: number }) =>
    api.post('/badges/equip', data),
}

export const fansubApi = {
  getAll: () => api.get('/fansubs'),
  
  getBySlug: (slug: string) => api.get(`/fansubs/${slug}`),
  
  create: (data: { name: string; slug: string; description: string }) =>
    api.post('/fansubs', data),
  
  join: (id: string) => api.post(`/fansubs/${id}/join`),
  
  approveMember: (id: string, data: { userId: string }) =>
    api.post(`/fansubs/${id}/approve`, data),
  
  leave: (id: string) => api.delete(`/fansubs/${id}/leave`),
}

export const shopApi = {
  getItems: () => api.get('/shop/items'),
  
  purchase: (data: { itemId: string; type: string; usePoints?: boolean }) =>
    api.post('/shop/purchase', data),
  
  getInventory: () => api.get('/shop/inventory'),
}

export const pointsApi = {
  getBalance: () => api.get('/points/balance'),
  
  earn: (data: { type: string; adId?: string }) => api.post('/points/earn', data),
  
  getHistory: () => api.get('/points/history'),
}

export const dmApi = {
  getConversations: () => api.get('/conversations'),
  
  createConversation: (data: { userId: string; type?: string }) =>
    api.post('/conversations', data),
  
  getMessages: (id: string, params?: { limit?: number; offset?: number }) =>
    api.get(`/conversations/${id}/messages`, { params }),
  
  sendMessage: (id: string, data: { content: string; type?: string }) =>
    api.post(`/conversations/${id}/messages`, data),
  
  deleteMessage: (id: string) => api.delete(`/messages/${id}`),
  
  markAsRead: (id: string) => api.post(`/messages/${id}/read`),
}

export const clipApi = {
  create: (data: { episodeId: string; startTime: number; endTime: number; title?: string }) =>
    api.post('/clips', data),
  
  getById: (id: string) => api.get(`/clips/${id}`),
  
  share: (id: string, data: { conversationId?: string }) =>
    api.post(`/clips/${id}/share`, data),
}

export const adminApi = {
  getDashboard: () => api.get('/admin/dashboard'),
  
  getUsers: () => api.get('/admin/users'),
  
  banUser: (id: string) => api.put(`/admin/users/${id}/ban`),
  
  getAllConversations: () => api.get('/admin/conversations/all'),
  
  getReports: () => api.get('/admin/reports'),
  
  resolveReport: (id: string, data: { status: string }) =>
    api.post(`/admin/reports/${id}/resolve`, data),
  
  approveVerification: (userId: string) =>
    api.post(`/admin/verification/${userId}/approve`),
  
  grantBadge: (data: { userId: string; badgeType: string; category: string }) =>
    api.post('/admin/badges/grant', data),
}