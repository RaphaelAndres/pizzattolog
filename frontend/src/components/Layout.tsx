import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import { LayoutDashboard, FileText, Bell, LogOut, ShieldCheck } from 'lucide-react'
import { useAuth } from '../hooks/useAuth'
import { useQuery } from '@tanstack/react-query'
import { dashboardService } from '../services/api'

export default function Layout() {
  const { usuario, logout } = useAuth()
  const navigate = useNavigate()

  const { data: alertas } = useQuery({
    queryKey: ['alertas'],
    queryFn: dashboardService.getAlertas,
    refetchInterval: 60000,
  })

  const alertCount = alertas?.total ?? 0

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const navClass = ({ isActive }: { isActive: boolean }) =>
    `flex items-center gap-3 px-4 py-2.5 rounded-lg text-sm font-medium transition-all ${
      isActive
        ? 'bg-orange-500 text-white shadow-sm'
        : 'text-slate-600 hover:bg-slate-100'
    }`

  return (
    <div className="flex h-screen bg-slate-50 font-sans">
      {/* Sidebar */}
      <aside className="w-64 bg-white border-r border-slate-200 flex flex-col">
        {/* Logo */}
        <div className="px-6 py-5 border-b border-slate-200">
          <div className="flex items-center gap-2">
            <ShieldCheck className="text-orange-500" size={24} />
            <div>
              <p className="font-bold text-slate-800 text-sm leading-none">PizzattoLog</p>
              <p className="text-xs text-slate-500 mt-0.5">Licenças</p>
            </div>
          </div>
        </div>

        {/* Nav */}
        <nav className="flex-1 px-3 py-4 space-y-1">
          <NavLink to="/dashboard" className={navClass}>
            <LayoutDashboard size={18} />
            Dashboard
          </NavLink>
          <NavLink to="/licencas" className={navClass}>
            <FileText size={18} />
            Licenças
          </NavLink>
          <NavLink to="/alertas" className={navClass}>
            <Bell size={18} />
            <span className="flex-1">Alertas</span>
            {alertCount > 0 && (
              <span className="bg-red-500 text-white text-xs font-bold px-1.5 py-0.5 rounded-full">
                {alertCount}
              </span>
            )}
          </NavLink>
        </nav>

        {/* User */}
        <div className="px-3 py-4 border-t border-slate-200">
          <div className="flex items-center gap-3 px-3 py-2">
            <div className="w-8 h-8 rounded-full bg-orange-100 flex items-center justify-center">
              <span className="text-orange-600 font-semibold text-sm">
                {usuario?.nome?.charAt(0).toUpperCase()}
              </span>
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-slate-800 truncate">{usuario?.nome}</p>
              <p className="text-xs text-slate-500 capitalize">{usuario?.role}</p>
            </div>
            <button
              onClick={handleLogout}
              className="text-slate-400 hover:text-red-500 transition-colors"
              title="Sair"
            >
              <LogOut size={16} />
            </button>
          </div>
        </div>
      </aside>

      {/* Content */}
      <main className="flex-1 overflow-auto">
        <Outlet />
      </main>
    </div>
  )
}
