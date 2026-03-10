import axios from 'axios'
import type {
  LoginRequest, LoginResponse, Licenca, ApiResponse,
  DashboardResumo, LicencaFiltros, Usuario,
} from '../types'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  timeout: 30000,
})

// Injeta token JWT automaticamente
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// Redireciona para login em caso de 401
api.interceptors.response.use(
  (r) => r,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('usuario')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// ── Auth ──────────────────────────────────────────────────────────────────────

export const authService = {
  login: (data: LoginRequest) =>
    api.post<LoginResponse>('/auth/login', data).then((r) => r.data),

  logout: () => api.post('/auth/logout'),
}

// ── Licenças ──────────────────────────────────────────────────────────────────

export const licencasService = {
  list: (filtros?: LicencaFiltros) =>
    api.get<ApiResponse<Licenca[]>>('/licencas', { params: filtros }).then((r) => r.data),

  getById: (id: number) =>
    api.get<ApiResponse<Licenca>>(`/licencas/${id}`).then((r) => r.data),

  create: (formData: FormData) =>
    api.post<ApiResponse<Licenca>>('/licencas', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }).then((r) => r.data),

  update: (id: number, formData: FormData) =>
    api.put<ApiResponse<Licenca>>(`/licencas/${id}`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }).then((r) => r.data),

  delete: (id: number) => api.delete(`/licencas/${id}`),

  getArquivoURL: (id: number) =>
    api.get<{ url: string }>(`/licencas/${id}/arquivo`).then((r) => r.data.url),
}

// ── Dashboard ─────────────────────────────────────────────────────────────────

export const dashboardService = {
  getResumo: () =>
    api.get<ApiResponse<DashboardResumo>>('/dashboard').then((r) => r.data),

  getAlertas: () =>
    api.get<ApiResponse<Licenca[]>>('/alertas').then((r) => r.data),
}

// ── Usuários ──────────────────────────────────────────────────────────────────

export const usuariosService = {
  list: () =>
    api.get<ApiResponse<Usuario[]>>('/usuarios').then((r) => r.data),

  create: (data: { nome: string; email: string; senha: string; role: string }) =>
    api.post<ApiResponse<Usuario>>('/usuarios', data).then((r) => r.data),

  update: (id: number, data: { nome: string; ativo: boolean }) =>
    api.put<ApiResponse<Usuario>>(`/usuarios/${id}`, data).then((r) => r.data),

  delete: (id: number) => api.delete(`/usuarios/${id}`),
}

export default api
