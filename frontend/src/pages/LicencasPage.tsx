import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { licencasService } from '../services/api'
import { format, differenceInDays } from 'date-fns'
import { Plus, Search, Filter, Trash2, Pencil, Download, FileText, Upload } from 'lucide-react'
import toast from 'react-hot-toast'
import type { Licenca, TipoLicenca } from '../types'
import { TIPO_LABELS, STATUS_LABELS, STATUS_COLORS } from '../types'

// ─── Modal Formulário ─────────────────────────────────────────────────────────

function ModalLicenca({
  licenca,
  onClose,
  onSalvo,
}: {
  licenca?: Licenca | null
  onClose: () => void
  onSalvo: () => void
}) {
  const [nome, setNome] = useState(licenca?.nome ?? '')
  const [tipo, setTipo] = useState<TipoLicenca>(licenca?.tipo ?? 'ambiental')
  const [orgao, setOrgao] = useState(licenca?.orgao_emissor ?? '')
  const [numero, setNumero] = useState(licenca?.numero ?? '')
  const [descricao, setDescricao] = useState(licenca?.descricao ?? '')
  const [dataVal, setDataVal] = useState(licenca?.data_validade?.substring(0, 10) ?? '')
  const [arquivo, setArquivo] = useState<File | null>(null)
  const [salvando, setSalvando] = useState(false)

  const salvar = async (e: React.FormEvent) => {
    e.preventDefault()
    setSalvando(true)
    const fd = new FormData()
    fd.append('nome', nome)
    fd.append('tipo', tipo)
    fd.append('orgao_emissor', orgao)
    fd.append('numero', numero)
    fd.append('descricao', descricao)
    fd.append('data_validade', dataVal)
    if (arquivo) fd.append('arquivo', arquivo)

    try {
      if (licenca) {
        await licencasService.update(licenca.id, fd)
        toast.success('Licença atualizada!')
      } else {
        await licencasService.create(fd)
        toast.success('Licença cadastrada!')
      }
      onSalvo()
    } catch {
      toast.error('Erro ao salvar licença')
    } finally {
      setSalvando(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
        <div className="px-6 py-5 border-b border-slate-200">
          <h2 className="font-bold text-slate-800">{licenca ? 'Editar Licença' : 'Nova Licença'}</h2>
        </div>
        <form onSubmit={salvar} className="px-6 py-5 space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="col-span-2">
              <label className="text-xs font-medium text-slate-600 mb-1 block">Nome da Licença *</label>
              <input value={nome} onChange={(e) => setNome(e.target.value)} required
                className="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500" />
            </div>
            <div>
              <label className="text-xs font-medium text-slate-600 mb-1 block">Tipo *</label>
              <select value={tipo} onChange={(e) => setTipo(e.target.value as TipoLicenca)} required
                className="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500">
                {(Object.entries(TIPO_LABELS) as [TipoLicenca, string][]).map(([v, l]) => (
                  <option key={v} value={v}>{l}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="text-xs font-medium text-slate-600 mb-1 block">Órgão Emissor</label>
              <input value={orgao} onChange={(e) => setOrgao(e.target.value)}
                className="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500" />
            </div>
            <div>
              <label className="text-xs font-medium text-slate-600 mb-1 block">Número</label>
              <input value={numero} onChange={(e) => setNumero(e.target.value)}
                className="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500" />
            </div>
            <div>
              <label className="text-xs font-medium text-slate-600 mb-1 block">Validade *</label>
              <input type="date" value={dataVal} onChange={(e) => setDataVal(e.target.value)} required
                className="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500" />
            </div>
            <div className="col-span-2">
              <label className="text-xs font-medium text-slate-600 mb-1 block">Descrição</label>
              <textarea value={descricao} onChange={(e) => setDescricao(e.target.value)} rows={2}
                className="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500 resize-none" />
            </div>
            <div className="col-span-2">
              <label className="text-xs font-medium text-slate-600 mb-1 block">Arquivo (PDF ou imagem)</label>
              <label className="flex items-center gap-2 border-2 border-dashed border-slate-200 rounded-lg px-4 py-3 cursor-pointer hover:border-orange-300 transition-colors">
                <Upload size={16} className="text-slate-400" />
                <span className="text-sm text-slate-500">
                  {arquivo ? arquivo.name : licenca?.arquivo_nome || 'Selecionar arquivo...'}
                </span>
                <input type="file" accept=".pdf,.png,.jpg,.jpeg" className="hidden"
                  onChange={(e) => setArquivo(e.target.files?.[0] ?? null)} />
              </label>
            </div>
          </div>

          <div className="flex justify-end gap-3 pt-2">
            <button type="button" onClick={onClose}
              className="px-4 py-2 text-sm text-slate-600 hover:bg-slate-100 rounded-lg transition-colors">
              Cancelar
            </button>
            <button type="submit" disabled={salvando}
              className="px-4 py-2 text-sm bg-orange-500 text-white font-medium rounded-lg hover:bg-orange-600 disabled:opacity-60 transition-colors">
              {salvando ? 'Salvando...' : 'Salvar'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

// ─── Página Principal ─────────────────────────────────────────────────────────

export default function LicencasPage() {
  const qc = useQueryClient()
  const [busca, setBusca] = useState('')
  const [filtroTipo, setFiltroTipo] = useState('')
  const [filtroStatus, setFiltroStatus] = useState('')
  const [modal, setModal] = useState<{ aberto: boolean; licenca?: Licenca | null }>({ aberto: false })

  const { data, isLoading } = useQuery({
    queryKey: ['licencas', busca, filtroTipo, filtroStatus],
    queryFn: () => licencasService.list({
      busca: busca || undefined,
      tipo: (filtroTipo as TipoLicenca) || undefined,
      status: (filtroStatus as Licenca['status']) || undefined,
    }),
  })

  const deleteMutation = useMutation({
    mutationFn: licencasService.delete,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['licencas'] })
      toast.success('Licença removida')
    },
    onError: () => toast.error('Erro ao remover'),
  })

  const handleDelete = (l: Licenca) => {
    if (confirm(`Remover "${l.nome}"?`)) deleteMutation.mutate(l.id)
  }

  const handleDownload = async (l: Licenca) => {
    if (!l.arquivo_key) return toast.error('Sem arquivo anexado')
    try {
      const url = await licencasService.getArquivoURL(l.id)
      window.open(url, '_blank')
    } catch {
      toast.error('Erro ao baixar arquivo')
    }
  }

  const fecharModal = () => setModal({ aberto: false })
  const salvoComSucesso = () => {
    qc.invalidateQueries({ queryKey: ['licencas'] })
    qc.invalidateQueries({ queryKey: ['dashboard'] })
    fecharModal()
  }

  const licencas = data?.dados ?? []

  return (
    <div className="p-8">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-slate-800">Licenças</h1>
          <p className="text-slate-500 text-sm mt-1">{data?.total ?? 0} licença(s) cadastrada(s)</p>
        </div>
        <button
          onClick={() => setModal({ aberto: true, licenca: null })}
          className="flex items-center gap-2 bg-orange-500 text-white text-sm font-medium px-4 py-2.5 rounded-lg hover:bg-orange-600 transition-colors"
        >
          <Plus size={16} />
          Nova Licença
        </button>
      </div>

      {/* Filtros */}
      <div className="flex flex-wrap gap-3 mb-6">
        <div className="relative flex-1 min-w-48">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <input
            value={busca}
            onChange={(e) => setBusca(e.target.value)}
            placeholder="Buscar por nome, número ou órgão..."
            className="w-full pl-9 pr-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"
          />
        </div>
        <select value={filtroTipo} onChange={(e) => setFiltroTipo(e.target.value)}
          className="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500">
          <option value="">Todos os tipos</option>
          {(Object.entries(TIPO_LABELS) as [TipoLicenca, string][]).map(([v, l]) => (
            <option key={v} value={v}>{l}</option>
          ))}
        </select>
        <select value={filtroStatus} onChange={(e) => setFiltroStatus(e.target.value)}
          className="border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500">
          <option value="">Todos os status</option>
          <option value="ativa">Ativa</option>
          <option value="proxima_vencimento">Próxima do Vencimento</option>
          <option value="vencida">Vencida</option>
        </select>
      </div>

      {/* Tabela */}
      {isLoading ? (
        <div className="flex justify-center py-16">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500" />
        </div>
      ) : licencas.length === 0 ? (
        <div className="text-center py-16 text-slate-400">
          <FileText size={48} className="mx-auto mb-3 opacity-30" />
          <p className="font-medium">Nenhuma licença encontrada</p>
          <p className="text-sm mt-1">Clique em "Nova Licença" para começar</p>
        </div>
      ) : (
        <div className="bg-white rounded-xl border border-slate-200 overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="bg-slate-50 border-b border-slate-200">
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wide">Licença</th>
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wide">Tipo</th>
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wide">Validade</th>
                <th className="text-left px-4 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wide">Status</th>
                <th className="px-4 py-3" />
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {licencas.map((l) => {
                const dias = differenceInDays(new Date(l.data_validade), new Date())
                return (
                  <tr key={l.id} className="hover:bg-slate-50 transition-colors">
                    <td className="px-4 py-3">
                      <p className="font-medium text-slate-800">{l.nome}</p>
                      {l.numero && <p className="text-xs text-slate-400">{l.numero}</p>}
                    </td>
                    <td className="px-4 py-3 text-slate-600">{TIPO_LABELS[l.tipo]}</td>
                    <td className="px-4 py-3">
                      <p className="text-slate-700">{format(new Date(l.data_validade), 'dd/MM/yyyy')}</p>
                      <p className={`text-xs ${dias < 0 ? 'text-red-500' : dias <= 30 ? 'text-amber-500' : 'text-slate-400'}`}>
                        {dias < 0 ? `Venceu há ${Math.abs(dias)} dias` : `${dias} dias restantes`}
                      </p>
                    </td>
                    <td className="px-4 py-3">
                      <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border ${STATUS_COLORS[l.status]}`}>
                        {STATUS_LABELS[l.status]}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-1 justify-end">
                        {l.arquivo_key && (
                          <button onClick={() => handleDownload(l)}
                            className="p-1.5 text-slate-400 hover:text-blue-500 hover:bg-blue-50 rounded-lg transition-colors" title="Baixar arquivo">
                            <Download size={14} />
                          </button>
                        )}
                        <button onClick={() => setModal({ aberto: true, licenca: l })}
                          className="p-1.5 text-slate-400 hover:text-orange-500 hover:bg-orange-50 rounded-lg transition-colors" title="Editar">
                          <Pencil size={14} />
                        </button>
                        <button onClick={() => handleDelete(l)}
                          className="p-1.5 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors" title="Remover">
                          <Trash2 size={14} />
                        </button>
                      </div>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}

      {modal.aberto && (
        <ModalLicenca licenca={modal.licenca} onClose={fecharModal} onSalvo={salvoComSucesso} />
      )}
    </div>
  )
}
