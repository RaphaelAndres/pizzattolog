import { useQuery } from '@tanstack/react-query'
import { dashboardService } from '../services/api'
import { differenceInDays, format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { AlertTriangle, Clock, XCircle } from 'lucide-react'
import { TIPO_LABELS } from '../types'

export default function AlertasPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['alertas'],
    queryFn: dashboardService.getAlertas,
    refetchInterval: 30000,
  })

  const licencas = data?.dados ?? []

  const criticas = licencas.filter((l) => {
    const dias = differenceInDays(new Date(l.data_validade), new Date())
    return dias <= 7 && dias >= 0
  })

  const atencao = licencas.filter((l) => {
    const dias = differenceInDays(new Date(l.data_validade), new Date())
    return dias > 7 && dias <= 30
  })

  const vencidas = licencas.filter((l) => l.status === 'vencida')

  if (isLoading) {
    return (
      <div className="p-8 flex justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500" />
      </div>
    )
  }

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-slate-800">Alertas</h1>
        <p className="text-slate-500 text-sm mt-1">
          Licenças que requerem atenção imediata
        </p>
      </div>

      {licencas.length === 0 && vencidas.length === 0 ? (
        <div className="bg-emerald-50 border border-emerald-200 rounded-xl p-8 text-center">
          <div className="text-5xl mb-3">✅</div>
          <p className="font-semibold text-emerald-800">Tudo em ordem!</p>
          <p className="text-emerald-600 text-sm mt-1">Nenhuma licença próxima do vencimento</p>
        </div>
      ) : (
        <div className="space-y-8">
          {/* Críticas (≤ 7 dias) */}
          {criticas.length > 0 && (
            <section>
              <div className="flex items-center gap-2 mb-4">
                <XCircle size={18} className="text-red-500" />
                <h2 className="font-semibold text-slate-800">Crítico — Vence em até 7 dias</h2>
                <span className="bg-red-100 text-red-700 text-xs font-bold px-2 py-0.5 rounded-full">{criticas.length}</span>
              </div>
              <div className="space-y-3">
                {criticas.map((l) => {
                  const dias = differenceInDays(new Date(l.data_validade), new Date())
                  return (
                    <div key={l.id} className="bg-white border border-red-200 rounded-xl p-4 flex items-center gap-4">
                      <div className="w-12 h-12 bg-red-100 rounded-xl flex items-center justify-center flex-shrink-0">
                        <Clock size={20} className="text-red-500" />
                      </div>
                      <div className="flex-1">
                        <p className="font-semibold text-slate-800">{l.nome}</p>
                        <p className="text-sm text-slate-500">{TIPO_LABELS[l.tipo]} · {l.orgao_emissor}</p>
                      </div>
                      <div className="text-right">
                        <p className="text-lg font-bold text-red-600">{dias}d</p>
                        <p className="text-xs text-slate-500">{format(new Date(l.data_validade), 'dd/MM/yyyy', { locale: ptBR })}</p>
                      </div>
                    </div>
                  )
                })}
              </div>
            </section>
          )}

          {/* Atenção (8–30 dias) */}
          {atencao.length > 0 && (
            <section>
              <div className="flex items-center gap-2 mb-4">
                <AlertTriangle size={18} className="text-amber-500" />
                <h2 className="font-semibold text-slate-800">Atenção — Vence em até 30 dias</h2>
                <span className="bg-amber-100 text-amber-700 text-xs font-bold px-2 py-0.5 rounded-full">{atencao.length}</span>
              </div>
              <div className="space-y-3">
                {atencao.map((l) => {
                  const dias = differenceInDays(new Date(l.data_validade), new Date())
                  return (
                    <div key={l.id} className="bg-white border border-amber-200 rounded-xl p-4 flex items-center gap-4">
                      <div className="w-12 h-12 bg-amber-100 rounded-xl flex items-center justify-center flex-shrink-0">
                        <AlertTriangle size={20} className="text-amber-500" />
                      </div>
                      <div className="flex-1">
                        <p className="font-semibold text-slate-800">{l.nome}</p>
                        <p className="text-sm text-slate-500">{TIPO_LABELS[l.tipo]} · {l.orgao_emissor}</p>
                      </div>
                      <div className="text-right">
                        <p className="text-lg font-bold text-amber-600">{dias}d</p>
                        <p className="text-xs text-slate-500">{format(new Date(l.data_validade), 'dd/MM/yyyy', { locale: ptBR })}</p>
                      </div>
                    </div>
                  )
                })}
              </div>
            </section>
          )}
        </div>
      )}
    </div>
  )
}
