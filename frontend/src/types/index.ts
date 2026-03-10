// ── Enums ─────────────────────────────────────────────────────────────────────

export type Role = 'admin' | 'gestor' | 'visualizador'

export type StatusLicenca = 'ativa' | 'proxima_vencimento' | 'vencida'

export type TipoLicenca =
  | 'ambiental'
  | 'policia_civil'
  | 'sanitaria'
  | 'bombeiros'
  | 'prefeitura'
  | 'outro'

// ── Entidades ────────────────────────────────────────────────────────────────

export interface Usuario {
  id: number
  nome: string
  email: string
  role: Role
  ativo: boolean
  criado_em: string
  atualizado_em: string
}

export interface Licenca {
  id: number
  nome: string
  tipo: TipoLicenca
  orgao_emissor: string
  numero: string
  descricao: string
  data_emissao: string | null
  data_validade: string
  status: StatusLicenca
  arquivo_key: string
  arquivo_nome: string
  arquivo_tamanho: number
  criado_por_id: number
  criado_por?: Usuario
  criado_em: string
  atualizado_em: string
  // Calculado no frontend
  dias_para_vencer?: number
}

// ── API ───────────────────────────────────────────────────────────────────────

export interface ApiResponse<T> {
  dados: T
  total?: number
}

export interface LoginRequest {
  email: string
  senha: string
}

export interface LoginResponse {
  token: string
  usuario: Usuario
}

export interface DashboardResumo {
  contagem_por_status: Record<StatusLicenca, number>
  proximas_vencer: Licenca[]
}

export interface LicencaFiltros {
  tipo?: TipoLicenca | ''
  status?: StatusLicenca | ''
  busca?: string
  ordem?: string
}

// ── UI ────────────────────────────────────────────────────────────────────────

export const TIPO_LABELS: Record<TipoLicenca, string> = {
  ambiental: 'Ambiental',
  policia_civil: 'Polícia Civil',
  sanitaria: 'Sanitária',
  bombeiros: 'Bombeiros',
  prefeitura: 'Prefeitura / Alvará',
  outro: 'Outro',
}

export const STATUS_LABELS: Record<StatusLicenca, string> = {
  ativa: 'Ativa',
  proxima_vencimento: 'Próxima do Vencimento',
  vencida: 'Vencida',
}

export const STATUS_COLORS: Record<StatusLicenca, string> = {
  ativa: 'text-emerald-700 bg-emerald-50 border-emerald-200',
  proxima_vencimento: 'text-amber-700 bg-amber-50 border-amber-200',
  vencida: 'text-red-700 bg-red-50 border-red-200',
}
