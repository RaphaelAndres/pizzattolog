import { createContext, useContext, useState, useCallback, ReactNode } from 'react'
import { authService } from '../services/api'
import type { Usuario } from '../types'

interface AuthContextType {
  usuario: Usuario | null
  token: string | null
  login: (email: string, senha: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [usuario, setUsuario] = useState<Usuario | null>(() => {
    const saved = localStorage.getItem('usuario')
    return saved ? JSON.parse(saved) : null
  })
  const [token, setToken] = useState<string | null>(() =>
    localStorage.getItem('token')
  )

  const login = useCallback(async (email: string, senha: string) => {
    const resp = await authService.login({ email, senha })
    localStorage.setItem('token', resp.token)
    localStorage.setItem('usuario', JSON.stringify(resp.usuario))
    setToken(resp.token)
    setUsuario(resp.usuario)
  }, [])

  const logout = useCallback(() => {
    authService.logout().catch(() => {})
    localStorage.removeItem('token')
    localStorage.removeItem('usuario')
    setToken(null)
    setUsuario(null)
  }, [])

  return (
    <AuthContext.Provider value={{ usuario, token, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth deve ser usado dentro de AuthProvider')
  return ctx
}
