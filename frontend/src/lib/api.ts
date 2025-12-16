// 1. AJUSTE DE URL: Apontando direto para o Go (Backend)
const API_BASE = 'http://localhost:8080/api'

export interface Company {
  id: number // Go usa number (int64), não string
  name: string
  // O Backend chama de contact_email, mas vamos mapear no Go para 'email' para facilitar
  email: string 
  cnpj: string // Adicionei CNPJ que faltava aqui
  // phone: string // Removi pois não criamos no Banco ainda
}

export interface User {
  id: number
  name: string
  email: string
  role: 'admin' | 'consultant'
}

export interface Contract {
  id: number
  companyId: number
  companyName?: string // Vem do JOIN
  title: string
  totalHours: number
  startDate: string
  endDate: string
  // Mudei de 'type' para 'contractType' para bater com o Go
  contractType: string 
  isActive: boolean
}

export interface Appointment {
  id: number
  contractId: number
  userId: number
  
  // A MUDANÇA CRÍTICA AQUI:
  // Removemos 'date' e 'hours' fixas
  // Adicionamos Start/End para o cálculo automático
  startTime: string 
  endTime: string
  
  description: string
  
  // Campos calculados/JOIN do Backend
  totalHours?: number 
  contractTitle?: string
  userName?: string
}

export interface ContractReport {
  contractId: string
  contractTitle: string
  companyName: string
  totalHours: number
  consumedHours: number
  remainingBalance: number
  idealMonthlyBurn: number
  actualBurnRate: number
  daysRemaining: number
  isOverBudget: boolean
}

// Helper genérico de Fetch (Mantive igual, só tipagem melhorada)
async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  })
  
  if (!res.ok) {
    // Tenta ler o erro do Go ou retorna genérico
    const errorData = await res.json().catch(() => null)
    throw new Error(errorData?.error || `Erro na API: ${res.statusText}`)
  }
  
  if (res.status === 204) {
    return {} as T
  }
  
  return res.json()
}

export const api = {
  companies: {
    list: () => fetchAPI<Company[]>('/companies'),
    get: (id: number) => fetchAPI<Company>(`/companies/${id}`),
    create: (data: any) => fetchAPI<Company>('/companies', { method: 'POST', body: JSON.stringify(data) }),
    update: (id: number, data: any) => fetchAPI<Company>(`/companies/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id: number) => fetchAPI<void>(`/companies/${id}`, { method: 'DELETE' }),
  },
  users: {
    list: () => fetchAPI<User[]>('/users'),
    get: (id: number) => fetchAPI<User>(`/users/${id}`),
    create: (data: any) => fetchAPI<User>('/users', { method: 'POST', body: JSON.stringify(data) }),
    update: (id: number, data: any) => fetchAPI<User>(`/users/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id: number) => fetchAPI<void>(`/users/${id}`, { method: 'DELETE' }),
  },
  contracts: {
    list: () => fetchAPI<Contract[]>('/contracts'),
    get: (id: number) => fetchAPI<Contract>(`/contracts/${id}`),
    create: (data: any) => fetchAPI<Contract>('/contracts', { method: 'POST', body: JSON.stringify(data) }),
    // Relatório de horas (Se não implementamos a rota /report ainda, isso vai dar 404)
    report: (id: number) => fetchAPI<any>(`/contracts/${id}/appointments`), 
  },
  appointments: {
    list: () => fetchAPI<Appointment[]>('/appointments'),
    // Agora o create exige startTime e endTime, não 'hours'
    create: (data: { contractId: number; userId: number; startTime: string; endTime: string; description: string }) =>
      fetchAPI<Appointment>('/appointments', { method: 'POST', body: JSON.stringify(data) }),
    delete: (id: number) => fetchAPI<void>(`/appointments/${id}`, { method: 'DELETE' }),
  },
}