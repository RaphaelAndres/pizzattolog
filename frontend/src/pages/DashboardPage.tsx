import { useQuery } from '@tanstack/react-query'
import { dashboardService } from '../services/api'
import { differenceInDays, format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { CheckCircle2, AlertTriangle, XCircle, Clock, FileText } from 'lucide-react'
import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer } from 'recharts'
import type { StatusLicenca } from '../types'

function StatCard({
  icon: Icon,
  label,
  value,
  color,
}: {
  icon: React.ElementType
  label: string
  value: number
  color: string
}) {
  return (
    <div className="bg-white rounded-xl border border-slate-200 p-5">
      <div className="flex items-center gap-3">
        <div className={`w-10 h-10 rounded-lg flex items-center justify-center ${color}`}>
          <Icon size={20} className="text-white" />
        </div>
        <div>
          <p className="text-2xl font-bold text-slate-800">{value}</p>
          <p className="text-xs text-slate-500 mt-0.5">{label}</p>
        </div>
      </div>
    </div>
  )
}

const PIE_COLORS = ['#10b981', '#f59e0b', '#ef4444']

export default function DashboardPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['dashboard'],
    queryFn: dashboardService.getResumo,
  })

  const resumo = data?.dados
  const contagem = resumo?.contagem_por_status ?? {} as Record<StatusLicenca, number>
  const proximas = resumo?.proximas_vencer ?? []

  const pieData = [
    { name: 'Ativas', value: contagem['ativa'] ?? 0 },
    { name: 'Próx. Vencimento', value: contagem['proxima_vencimento'] ?? 0 },
    { name: 'Vencidas', value: contagem['vencida'] ?? 0 },
  ]

  if (isLoading) {
    return (
      <div className="p-8 flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500" />
      </div>
    )
  }

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-slate-800">Dashboard</h1>
        <p className="text-slate-500 text-sm mt-1">Visão geral das licenças operacionais</p>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
        <StatCard
          icon={CheckCircle2}
          label="Licenças Ativas"
          value={contagem['ativa'] ?? 0}
          color="bg-emerald-500"
        />
        <StatCard
          icon={AlertTriangle}
          label="Próximas do Vencimento"
          value={contagem['proxima_vencimento'] ?? 0}
          color="bg-amber-500"
        />
        <StatCard
          icon={XCircle}
          label="Vencidas"
          value={contagem['vencida'] ?? 0}
          color="bg-red-500"
        />
      </div>

      {/* Chart + Proximas */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Gráfico de pizza */}
        <div className="bg-white rounded-xl border border-slate-200 p-6">
          <h2 className="font-semibold text-slate-700 mb-4 text-sm">Distribuição por Status</h2>
          <ResponsiveContainer width="100%" height={220}>
            <PieChart>
              <Pie data={pieData} cx="50%" cy="50%" innerRadius={60} outerRadius={90} dataKey="value" paddingAngle={3}>
                {pieData.map((_, i) => (
                  <Cell key={i} fill={PIE_COLORS[i]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
          <div className="flex justify-center gap-4 mt-2">
            {pieData.map((item, i) => (
              <div key={i} className="flex items-center gap-1.5">
                <div className="w-2.5 h-2.5 rounded-full" style={{ background: PIE_COLORS[i] }} />
                <span className="text-xs text-slate-600">{item.name}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Próximas a vencer */}
        <div className="bg-white rounded-xl border border-slate-200 p-6">
          <h2 className="font-semibold text-slate-700 mb-4 text-sm">Próximas a Vencer</h2>
          {proximas.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-32 text-slate-400">
              <FileText size={32} className="mb-2 opacity-40" />
              <p className="text-sm">Nenhuma licença próxima do vencimento</p>
            </div>
          ) : (
            <div className="space-y-3">
              {proximas.map((l) => {
                const dias = differenceInDays(new Date(l.data_validade), new Date())
                return (
                  <div key={l.id} className="flex items-center gap-3 p-3 rounded-lg bg-slate-50 hover:bg-slate-100 transition-colors">
                    <div className="flex-shrink-0">
                      <Clock size={16} className={dias <= 7 ? 'text-red-500' : 'text-amber-500'} />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-slate-700 truncate">{l.nome}</p>
                      <p className="text-xs text-slate-500">
                        Vence em {format(new Date(l.data_validade), 'dd/MM/yyyy', { locale: ptBR })}
                      </p>
                    </div>
                    <span className={`text-xs font-semibold px-2 py-1 rounded-lg ${dias <= 7 ? 'bg-red-100 text-red-700' : 'bg-amber-100 text-amber-700'}`}>
                      {dias}d
                    </span>
                  </div>
                )
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
