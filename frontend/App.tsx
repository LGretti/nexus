import React, { useState, useEffect } from 'react';
import { HashRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom';
import { 
  LayoutDashboard, 
  Building2, 
  FileText, 
  Users, 
  Clock, 
  Plus, 
  Trash2, 
  Edit2, 
  AlertCircle, 
  Loader2,
  X,
  BarChart3,
  Calendar,
  PieChart,
  ArrowRight,
  Ghost,
  MonitorOff
} from 'lucide-react';
import * as api from './services/api';
import { Company, Contract, User, Appointment } from './types';

// --- Utility Components ---

const Alert = ({ message, type = 'error' }: { message: string | null, type?: 'error' | 'success' }) => {
  if (!message) return null;
  const bg = type === 'error' ? 'bg-red-50' : 'bg-green-50';
  const text = type === 'error' ? 'text-red-700' : 'text-green-700';
  const border = type === 'error' ? 'border-red-200' : 'border-green-200';
  return (
    <div className={`p-4 mb-4 rounded-md border ${bg} ${border} ${text} flex items-center gap-2 animate-in fade-in slide-in-from-top-2`}>
      <AlertCircle className="w-5 h-5" />
      <span>{message}</span>
    </div>
  );
};

export const formatDuration = (totalSeconds: number): string => {
  if (!totalSeconds) return "0h00";

  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);

  // padStart(2, '0') garante o "05" em vez de só "5"
  return `${hours}h${minutes.toString().padStart(2, '0')}`;
};

const Button = ({ children, onClick, variant = 'primary', className = '', disabled = false, type = "button" }: any) => {
  const baseStyle = "inline-flex items-center justify-center px-4 py-2 rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed text-sm";
  const variants = {
    primary: "bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500",
    secondary: "bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 focus:ring-gray-500",
    danger: "bg-red-600 text-white hover:bg-red-700 focus:ring-red-500",
    ghost: "text-gray-600 hover:bg-gray-100",
  };
  
  return (
    <button 
      type={type}
      className={`${baseStyle} ${variants[variant as keyof typeof variants]} ${className}`} 
      onClick={onClick} 
      disabled={disabled}
    >
      {children}
    </button>
  );
};

const EmptyState = ({ title, description, actionLabel, onAction }: any) => (
  <div className="bg-white rounded-xl border border-dashed border-gray-300 p-12 text-center flex flex-col items-center justify-center animate-in fade-in zoom-in-95 duration-200">
    <div className="bg-blue-50 p-4 rounded-full mb-4">
      <Ghost className="w-8 h-8 text-blue-500" />
    </div>
    <h3 className="text-lg font-semibold text-gray-900 mb-2">{title}</h3>
    <p className="text-gray-500 max-w-sm mb-6">{description}</p>
    {actionLabel && onAction && (
      <Button onClick={onAction}>
        <Plus className="w-4 h-4 mr-2" />
        {actionLabel}
      </Button>
    )}
  </div>
);

const Modal = ({ isOpen, onClose, title, children, size = 'md' }: any) => {
  if (!isOpen) return null;
  const sizeClasses = {
    md: 'max-w-lg',
    lg: 'max-w-4xl'
  };
  
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm overflow-y-auto">
      <div className={`bg-white rounded-lg shadow-xl w-full ${sizeClasses[size as keyof typeof sizeClasses]} overflow-hidden animate-in zoom-in-95 duration-200`}>
        <div className="px-6 py-4 border-b flex justify-between items-center bg-gray-50">
          <h3 className="text-lg font-semibold text-gray-900">{title}</h3>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-500 transition-colors">
            <X className="w-5 h-5" />
          </button>
        </div>
        <div className="p-6">{children}</div>
      </div>
    </div>
  );
};

const Table = ({ headers, children }: { headers: string[], children?: React.ReactNode }) => (
  <div className="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
    <div className="overflow-x-auto">
      <table className="w-full text-sm text-left">
        <thead className="bg-gray-50 text-gray-500 font-medium border-b border-gray-200">
          <tr>
            {headers.map((h, i) => (
              <th key={i} className="px-6 py-3 whitespace-nowrap">{h}</th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100 text-gray-700">
          {children}
        </tbody>
      </table>
    </div>
  </div>
);

const formatDate = (dateString: string) => {
  if (!dateString) return '-';
  return new Date(dateString).toLocaleString('pt-BR', {
    day: '2-digit', month: '2-digit', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  });
};

const formatCurrency = (val: number) => val.toFixed(1);

// --- Dashboard Component ---

const Dashboard = () => {
  const [stats, setStats] = useState({
    activeContracts: 0,
    totalCompanies: 0,
    recentActivity: [] as Appointment[]
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [companies, contracts, appointments] = await Promise.all([
          api.getCompanies(),
          api.getContracts(),
          api.getAppointments()
        ]);

        const safeContracts = Array.isArray(contracts) ? contracts : [];
        const safeCompanies = Array.isArray(companies) ? companies : [];
        const safeAppointments = Array.isArray(appointments) ? appointments : [];

        const activeContracts = safeContracts.filter(c => c.isActive).length;
        // Sort appointments by date descending
        const recent = [...safeAppointments]
          .sort((a, b) => new Date(b.startTime).getTime() - new Date(a.startTime).getTime())
          .slice(0, 5);

        setStats({
          activeContracts,
          totalCompanies: safeCompanies.length,
          recentActivity: recent
        });
      } catch (e: any) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const Card = ({ title, value, icon: Icon, color }: any) => (
    <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm flex items-center gap-4">
      <div className={`p-3 rounded-lg ${color}`}>
        <Icon className="w-6 h-6 text-white" />
      </div>
      <div>
        <p className="text-sm font-medium text-gray-500">{title}</p>
        <p className="text-2xl font-bold text-gray-900">{loading ? '...' : value}</p>
      </div>
    </div>
  );

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Visão Geral</h1>
      
      <Alert message={error} />

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card title="Active Contracts" value={stats.activeContracts} icon={FileText} color="bg-indigo-500" />
        <Card title="Total Companies" value={stats.totalCompanies} icon={Building2} color="bg-blue-500" />
        <Card title="Recent Activities" value={stats.recentActivity.length} icon={Clock} color="bg-emerald-500" />
      </div>
      
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Activity List */}
        <div className="bg-white rounded-xl border border-gray-200 shadow-sm flex flex-col">
          <div className="p-6 border-b border-gray-100">
            <h3 className="text-lg font-semibold text-gray-900 flex items-center gap-2">
              <Clock className="w-5 h-5 text-gray-400" /> Atividade Recente
            </h3>
          </div>
          <div className="p-0 flex-1">
            {stats.recentActivity.length === 0 ? (
               <div className="p-6 text-center text-gray-500">Nenhuma atividade recente encontrada.</div>
            ) : (
              <div className="divide-y divide-gray-100">
                {stats.recentActivity.map((apt) => (
                  <div key={apt.id} className="p-4 hover:bg-gray-50 transition-colors flex items-start gap-3">
                    <div className="mt-1">
                      <div className="w-2 h-2 rounded-full bg-blue-500"></div>
                    </div>
                    <div className="flex-1">
                      <p className="text-sm font-medium text-gray-900">
                        {apt.userName} <span className="text-gray-500 font-normal">usou</span> {formatDuration(apt.durationSeconds)}
                      </p>
                      <p className="text-xs text-blue-600 mt-0.5">{apt.contractTitle}</p>
                      <p className="text-xs text-gray-500 mt-1">{apt.description}</p>
                    </div>
                    <div className="text-xs text-gray-400 whitespace-nowrap">
                      {new Date(apt.startTime).toLocaleDateString()}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
          <div className="p-4 border-t border-gray-100 bg-gray-50 rounded-b-xl">
             <Link to="/appointments" className="text-sm text-blue-600 font-medium hover:text-blue-800 flex items-center gap-1">
               Ver todos os atendimentos <ArrowRight className="w-4 h-4" />
             </Link>
          </div>
        </div>

        {/* Welcome / Quick Actions */}
        <div className="bg-gradient-to-br from-slate-900 to-slate-800 text-white rounded-xl shadow-sm p-8 flex flex-col justify-center">
          <h2 className="text-2xl font-bold mb-2">Bem-vindo de volta</h2>
          <p className="text-slate-300 mb-6">
            Gerencie seus contratos e acompanhe o tempo de forma eficiente. Você tem {stats.activeContracts} contratos ativos em andamento.
          </p>
          <div className="flex flex-wrap gap-3">
             <Link to="/contracts" className="bg-white text-slate-900 px-4 py-2 rounded-lg font-medium hover:bg-gray-100 transition-colors">
               Gerenciar Contratos
             </Link>
             <Link to="/appointments" className="bg-slate-700 text-white px-4 py-2 rounded-lg font-medium hover:bg-slate-600 transition-colors">
               Registrar Atividade
             </Link>
          </div>
        </div>
      </div>
    </div>
  );
};

// --- Contracts Page with Report Modal ---

const ContractReportModal = ({ contract, isOpen, onClose }: { contract: Contract | null, isOpen: boolean, onClose: () => void }) => {
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (isOpen && contract) {
      setLoading(true);
      api.getContractAppointments(contract.id)
        .then((data) => setAppointments(data || [])) // Safety check
        .catch(console.error)
        .finally(() => setLoading(false));
    }
  }, [isOpen, contract]);

  if (!contract) return null;

  // Calculations
  const totalHours = contract.totalHours;
  const consumedHours = appointments.reduce((sum, a) => sum + ((a.durationSeconds || 0) / 3600), 0);
  const remainingHours = totalHours - consumedHours;
  const progress = Math.min((consumedHours / totalHours) * 100, 100);

  // Date Logic
  const start = new Date(contract.startDate);
  const end = new Date(contract.endDate);
  const now = new Date();
  const totalDurationMs = end.getTime() - start.getTime();
  const elapsedMs = now.getTime() - start.getTime();
  const daysRemaining = Math.max(0, Math.ceil((end.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)));
  
  // Burn Rate (Hours per month)
  // Approximate months elapsed. Avoid division by zero.
  const monthsElapsed = Math.max(elapsedMs / (1000 * 60 * 60 * 24 * 30.44), 0.5); 
  const actualBurnRate = consumedHours / monthsElapsed;
  const totalMonths = totalDurationMs / (1000 * 60 * 60 * 24 * 30.44);
  const idealBurnRate = totalHours / totalMonths;

  const Stat = ({ label, value, subtext, color = "text-gray-900" }: any) => (
    <div className="bg-gray-50 p-4 rounded-lg border border-gray-100">
      <p className="text-xs font-medium text-gray-500 uppercase">{label}</p>
      <p className={`text-2xl font-bold mt-1 ${color}`}>{value}</p>
      {subtext && <p className="text-xs text-gray-400 mt-1">{subtext}</p>}
    </div>
  );

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Detalhes do Contrato" size="lg">
      <div className="space-y-6">
        {/* Header Info */}
        <div className="flex justify-between items-start">
          <div>
            <h2 className="text-xl font-bold text-gray-900">{contract.title}</h2>
            <p className="text-gray-500">{contract.companyName}</p>
          </div>
          <div className="text-right">
             <span className={`px-2 py-1 rounded text-xs font-medium ${contract.isActive ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                {contract.isActive ? 'Active' : 'Inactive'}
             </span>
             <p className="text-xs text-gray-400 mt-1">{new Date(contract.startDate).toLocaleDateString()} - {new Date(contract.endDate).toLocaleDateString()}</p>
          </div>
        </div>

        {/* Progress Bar */}
        <div className="space-y-2">
          <div className="flex justify-between text-sm">
            <span className="font-medium text-gray-700">Tempo Utilizado</span>
            <span className="text-gray-500">{consumedHours.toFixed(1)} / {totalHours} hrs</span>
          </div>
          <div className="h-3 bg-gray-100 rounded-full overflow-hidden">
            <div 
              className={`h-full ${remainingHours < 0 ? 'bg-red-500' : 'bg-blue-600'}`} 
              style={{ width: `${progress}%` }}
            />
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <Stat label="Horas Totais" value={`${totalHours}h`} />
          <Stat label="Consumidas" value={`${consumedHours.toFixed(1)}h`} />
          <Stat 
            label="Saldo" 
            value={`${remainingHours.toFixed(1)}h`} 
            color={remainingHours < 0 ? 'text-red-600' : 'text-green-600'} 
          />
          <Stat label="Tempo Restante" value={`${daysRemaining} dias`} />
        </div>

        {/* Burn Rate Analysis */}
        <div className="bg-indigo-50 p-4 rounded-lg border border-indigo-100 flex items-center justify-between">
           <div>
             <h4 className="font-semibold text-indigo-900">Média de Consumo Ideal</h4>
             <p className="text-sm text-indigo-700">
               Atual: <b>{actualBurnRate.toFixed(1)} h/mês</b> vs Ideal: <b>{idealBurnRate.toFixed(1)} h/mês</b>
             </p>
           </div>
           {actualBurnRate > idealBurnRate * 1.2 && (
             <div className="bg-red-100 text-red-700 px-3 py-1 rounded text-xs font-bold">
               High Usage
             </div>
           )}
        </div>

        {/* Table */}
        <div className="pt-4 border-t border-gray-100">
          <h4 className="font-semibold text-gray-900 mb-4">Atendimentos Realizados</h4>
          {loading ? (
            <div className="text-center py-8"><Loader2 className="w-6 h-6 animate-spin mx-auto text-blue-600"/></div>
          ) : (
            <div className="max-h-60 overflow-y-auto border rounded-lg">
              <table className="w-full text-sm text-left">
                <thead className="bg-gray-50 sticky top-0">
                  <tr>
                    <th className="px-4 py-2">Data</th>
                    <th className="px-4 py-2">Consultor</th>
                    <th className="px-4 py-2">Descrição</th>
                    <th className="px-4 py-2 text-right">Horas</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {appointments.map(a => (
                    <tr key={a.id}>
                      <td className="px-4 py-2 whitespace-nowrap text-gray-500">{new Date(a.startTime).toLocaleDateString()}</td>
                      <td className="px-4 py-2">{a.userName}</td>
                      <td className="px-4 py-2 text-gray-600 truncate max-w-xs">{a.description}</td>
                      <td className="px-4 py-2 text-right font-medium">{formatDuration(a.durationSeconds || 0)}</td>
                    </tr>
                  ))}
                  {appointments.length === 0 && (
                    <tr><td colSpan={4} className="p-4 text-center text-gray-400">Ainda não foram registradas atividades nesse contrato.</td></tr>
                  )}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </Modal>
  );
};

const ContractsPage = () => {
  const [contracts, setContracts] = useState<Contract[]>([]);
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Modal States
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [selectedContract, setSelectedContract] = useState<Contract | null>(null);

  const [formData, setFormData] = useState({
    companyId: 0,
    title: '',
    contractType: 'Mensal',
    totalHours: 0,
    startDate: '',
    endDate: '',
    isActive: true
  });

  const loadData = async () => {
    try {
      setLoading(true);
      const [cData, compData] = await Promise.all([api.getContracts(), api.getCompanies()]);
      setContracts(cData || []);
      setCompanies(compData || []);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { loadData(); }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      const payload = {
        ...formData,
        startDate: new Date(formData.startDate).toISOString(),
        endDate: new Date(formData.endDate).toISOString(),
        companyId: Number(formData.companyId),
        totalHours: Number(formData.totalHours)
      };
      
      await api.createContract(payload);
      setIsCreateOpen(false);
      setFormData({
        companyId: 0,
        title: '',
        contractType: 'Mensal',
        totalHours: 0,
        startDate: '',
        endDate: '',
        isActive: true
      });
      loadData();
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold text-gray-900">Contratos</h1>
        <Button onClick={() => setIsCreateOpen(true)}><Plus className="w-4 h-4 mr-2" /> Adicionar Contrato</Button>
      </div>

      <Alert message={error} />

      {loading ? (
        <Loader2 className="w-8 h-8 animate-spin text-blue-600 mx-auto" />
      ) : contracts.length === 0 ? (
        <EmptyState 
          title="Nenhum Contrato Encontrado" 
          description="Crie seu primeiro contrato para começar a acompanhar horas e atividades."
          actionLabel="Criar Contrato"
          onAction={() => setIsCreateOpen(true)}
        />
      ) : (
        <Table headers={['Título', 'Empresa', 'Tipo', 'Horas', 'Status', 'Duração', 'Ações']}>
          {contracts.map(c => (
            <tr key={c.id} className="hover:bg-gray-50 transition-colors">
              <td className="px-6 py-4 font-medium text-gray-900">{c.title}</td>
              <td className="px-6 py-4 text-blue-600">{c.companyName}</td>
              <td className="px-6 py-4"><span className="px-2 py-1 bg-gray-100 rounded text-xs text-gray-600">{c.contractType}</span></td>
              <td className="px-6 py-4 font-mono">{c.totalHours}h</td>
              <td className="px-6 py-4">
                <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${c.isActive ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                  {c.isActive ? 'Ativo' : 'Inativo'}
                </span>
              </td>
              <td className="px-6 py-4 text-xs text-gray-500">
                {new Date(c.startDate).toLocaleDateString()} - {new Date(c.endDate).toLocaleDateString()}
              </td>
              <td className="px-6 py-4">
                <Button variant="ghost" onClick={() => setSelectedContract(c)} className="text-xs px-2 py-1">
                   <BarChart3 className="w-4 h-4 mr-1"/> Relatório
                </Button>
              </td>
            </tr>
          ))}
        </Table>
      )}

      {/* Report Modal */}
      <ContractReportModal 
        contract={selectedContract} 
        isOpen={!!selectedContract} 
        onClose={() => setSelectedContract(null)} 
      />

      {/* Create Modal */}
      <Modal isOpen={isCreateOpen} onClose={() => setIsCreateOpen(false)} title="Novo Contrato">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Empresa</label>
            <select required className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formData.companyId} onChange={e => setFormData({...formData, companyId: Number(e.target.value)})}>
              <option value={0}>Selecione uma empresa</option>
              {companies.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Título</label>
            <input required type="text" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formData.title} onChange={e => setFormData({...formData, title: e.target.value})} />
          </div>
          <div className="grid grid-cols-2 gap-4">
             <div>
              <label className="block text-sm font-medium text-gray-700">Tipo</label>
              <input required type="text" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formData.contractType} onChange={e => setFormData({...formData, contractType: e.target.value})} />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Total de Horas</label>
              <input required type="number" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formatDuration(formData.durationSeconds)} onChange={e => setFormData({...formData, totalHours: Number(e.target.value)})} />
            </div>
          </div>
          <div className="grid grid-cols-2 gap-4">
             <div>
              <label className="block text-sm font-medium text-gray-700">Data de Início</label>
              <input required type="date" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formData.startDate} onChange={e => setFormData({...formData, startDate: e.target.value})} />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Data de Término</label>
              <input required type="date" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formData.endDate} onChange={e => setFormData({...formData, endDate: e.target.value})} />
            </div>
          </div>
          <div className="flex items-center gap-2 mt-2">
            <input type="checkbox" id="isActive" checked={formData.isActive} onChange={e => setFormData({...formData, isActive: e.target.checked})} className="rounded text-blue-600 focus:ring-blue-500 h-4 w-4" />
            <label htmlFor="isActive" className="text-sm text-gray-700">O contrato está Ativo</label>
          </div>
          <div className="flex justify-end gap-3 mt-6">
            <Button variant="secondary" onClick={() => setIsCreateOpen(false)}>Cancelar</Button>
            <Button type="submit">Adicionar Contrato</Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

// --- Other Pages (Company, User, Appointments) ---

const CompaniesPage = () => {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [formData, setFormData] = useState({ name: '', cnpj: '', email: '' });

  const loadData = async () => {
    try {
      setLoading(true);
      const data = await api.getCompanies();
      setCompanies(data || []);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { loadData(); }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      if (editingId) {
        await api.updateCompany(editingId, formData);
      } else {
        await api.createCompany(formData);
      }
      setIsModalOpen(false);
      setFormData({ name: '', cnpj: '', email: '' });
      setEditingId(null);
      loadData();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleEdit = (c: Company) => {
    setFormData({ name: c.name, cnpj: c.cnpj, email: c.email });
    setEditingId(c.id);
    setIsModalOpen(true);
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm("Você deseja desativar esta empresa?")) return;
    try {
      await api.deleteCompany(id);
      loadData();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const openCreate = () => {
    setEditingId(null);
    setFormData({ name: '', cnpj: '', email: '' });
    setIsModalOpen(true);
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold text-gray-900">Empresas</h1>
        <Button onClick={openCreate}><Plus className="w-4 h-4 mr-2" /> Adicionar Empresa</Button>
      </div>

      <Alert message={error} />

      {loading ? (
        <Loader2 className="w-8 h-8 animate-spin text-blue-600 mx-auto" /> 
      ) : companies.length === 0 ? (
        <EmptyState 
          title="No Companies Found" 
          description="Get started by adding your first client company."
          actionLabel="Add Company"
          onAction={openCreate}
        />
      ) : (
        <Table headers={['Nome', 'CNPJ', 'Email', 'Ações']}>
          {companies.map(c => (
            <tr key={c.id} className="hover:bg-gray-50 transition-colors">
              <td className="px-6 py-4 font-medium text-gray-900">{c.name}</td>
              <td className="px-6 py-4 text-gray-500 font-mono text-xs">{c.cnpj}</td>
              <td className="px-6 py-4">{c.email}</td>
              <td className="px-6 py-4">
                <div className="flex gap-2">
                  <button onClick={() => handleEdit(c)} className="p-1 hover:bg-blue-100 text-blue-600 rounded"><Edit2 className="w-4 h-4" /></button>
                  <button onClick={() => handleDelete(c.id)} className="p-1 hover:bg-red-100 text-red-600 rounded"><MonitorOff className="w-4 h-4" /></button>
                </div>
              </td>
            </tr>
          ))}
        </Table>
      )}

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={editingId ? "Editar Empresa" : "Nova Empresa"}>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Nome</label>
            <input required type="text" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2 focus:ring-blue-500 focus:border-blue-500" value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">CNPJ</label>
            <input required type="text" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2 focus:ring-blue-500 focus:border-blue-500" value={formData.cnpj} onChange={e => setFormData({...formData, cnpj: e.target.value})} placeholder="00.000.000/0000-00" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">E-mail</label>
            <input required type="email" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2 focus:ring-blue-500 focus:border-blue-500" value={formData.email} onChange={e => setFormData({...formData, email: e.target.value})} />
          </div>
          <div className="flex justify-end gap-3 mt-6">
            <Button variant="secondary" onClick={() => setIsModalOpen(false)}>Cancelar</Button>
            <Button type="submit">{editingId ? 'Salvar Alterações' : 'Adicionar Empresa'}</Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

const UsersPage = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [formData, setFormData] = useState<Omit<User, 'id'>>({ name: '', email: '', role: 'consultant' });

  const loadData = async () => {
    try {
      setLoading(true);
      const data = await api.getUsers();
      setUsers(data || []);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { loadData(); }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.createUser(formData);
      setIsModalOpen(false);
      setFormData({ name: '', email: '', role: 'consultant' });
      loadData();
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold text-gray-900">Usuários</h1>
        <Button onClick={() => setIsModalOpen(true)}><Plus className="w-4 h-4 mr-2" /> Adicionar Usuário</Button>
      </div>
      
      <Alert message={error} />

      {loading ? (
        <Loader2 className="w-8 h-8 animate-spin text-blue-600 mx-auto" />
      ) : users.length === 0 ? (
        <EmptyState 
          title="Nenhum Usuário Encontrado" 
          description="Adicione consultores ou administradores ao sistema."
          actionLabel="Adicionar Usuário"
          onAction={() => setIsModalOpen(true)}
        />
      ) : (
        <Table headers={['Nome', 'E-mail', 'Função']}>
          {users.map(u => (
            <tr key={u.id} className="hover:bg-gray-50">
              <td className="px-6 py-4 font-medium text-gray-900">{u.name}</td>
              <td className="px-6 py-4 text-gray-500">{u.email}</td>
              <td className="px-6 py-4">
                <span className={`px-2 py-1 rounded-full text-xs font-semibold ${u.role === 'admin' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'}`}>
                  {u.role.toUpperCase()}
                </span>
              </td>
            </tr>
          ))}
        </Table>
      )}

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title="New User">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Nome</label>
            <input required type="text" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">E-mail</label>
            <input required type="email" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formData.email} onChange={e => setFormData({...formData, email: e.target.value})} />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Função</label>
            <select className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formData.role} onChange={e => setFormData({...formData, role: e.target.value as any})}>
              <option value="consultant">Consultor</option>
              <option value="admin">Administrador</option>
            </select>
          </div>
          <div className="flex justify-end gap-3 mt-6">
            <Button variant="secondary" onClick={() => setIsModalOpen(false)}>Cancelar</Button>
            <Button type="submit">Adicionar Usuário</Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

const AppointmentsPage = () => {
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [contracts, setContracts] = useState<Contract[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  
  const [formData, setFormData] = useState({
    contractId: 0,
    userId: 0,
    startTime: '',
    endTime: '',
    description: ''
  });

  const loadData = async () => {
    try {
      setLoading(true);
      const [appData, contData, usrData] = await Promise.all([
        api.getAppointments(),
        api.getContracts(),
        api.getUsers()
      ]);
      setAppointments(appData || []);
      setContracts(contData || []);
      setUsers(usrData || []);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { loadData(); }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      // Ensure ISO conversion handles Timezone offset if browser is not UTC
      const payload = {
        ...formData,
        contractId: Number(formData.contractId),
        userId: Number(formData.userId),
        startTime: new Date(formData.startTime).toISOString(),
        endTime: new Date(formData.endTime).toISOString(),
      };
      
      await api.createAppointment(payload);
      setIsModalOpen(false);
      setFormData({ contractId: 0, userId: 0, startTime: '', endTime: '', description: '' });
      loadData();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm("Deseja excluir esse registro?")) return;
    try {
      await api.deleteAppointment(id);
      loadData();
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold text-gray-900">Atendimentos</h1>
        <Button onClick={() => setIsModalOpen(true)}><Plus className="w-4 h-4 mr-2" /> Registrar Horas</Button>
      </div>

      <Alert message={error} />

      {loading ? (
        <Loader2 className="w-8 h-8 animate-spin text-blue-600 mx-auto" />
      ) : appointments.length === 0 ? (
        <EmptyState 
          title="Nenhum Atendimento Registrado" 
          description="Registre horas de trabalho para seus contratos ativos."
          actionLabel="Registrar Horas"
          onAction={() => setIsModalOpen(true)}
        />
      ) : (
        <Table headers={['Data', 'Usuário', 'Contrato', 'Descrição', 'Horas', 'Ação']}>
          {appointments.map(a => (
            <tr key={a.id} className="hover:bg-gray-50">
              <td className="px-6 py-4 text-sm text-gray-600">
                <div className="font-medium text-gray-900">{formatDate(a.startTime)}</div>
                <div className="text-xs">até {formatDate(a.endTime).split(' ')[2]}</div>
              </td>
              <td className="px-6 py-4 font-medium text-gray-900">{a.userName}</td>
              <td className="px-6 py-4 text-blue-600 text-sm">{a.contractTitle}</td>
              <td className="px-6 py-4 text-gray-500 truncate max-w-xs">{a.description}</td>
              <td className="px-6 py-4 font-bold text-gray-800">{formatDuration(a.durationSeconds)}</td>
              <td className="px-6 py-4">
                <button onClick={() => handleDelete(a.id)} className="p-1 hover:bg-red-100 text-red-600 rounded"><Trash2 className="w-4 h-4" /></button>
              </td>
            </tr>
          ))}
        </Table>
      )}

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title="Log Time">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Contrato</label>
            <select required className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formData.contractId} onChange={e => setFormData({...formData, contractId: Number(e.target.value)})}>
              <option value={0}>Selecione o Contrato</option>
              {contracts.filter(c => c.isActive).map(c => <option key={c.id} value={c.id}>{c.title}</option>)}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Consultor</label>
            <select required className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formData.userId} onChange={e => setFormData({...formData, userId: Number(e.target.value)})}>
              <option value={0}>Selecione o Consultor</option>
              {users.map(u => <option key={u.id} value={u.id}>{u.name}</option>)}
            </select>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">Hora de Início</label>
              <input required type="datetime-local" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formData.startTime} onChange={e => setFormData({...formData, startTime: e.target.value})} />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Hora de Término</label>
              <input required type="datetime-local" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formData.endTime} onChange={e => setFormData({...formData, endTime: e.target.value})} />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Descrição</label>
            <textarea required rows={3} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm border p-2" value={formData.description} onChange={e => setFormData({...formData, description: e.target.value})} placeholder="Descreva o atendimento realizado" />
          </div>
          <div className="flex justify-end gap-3 mt-6">
            <Button variant="secondary" onClick={() => setIsModalOpen(false)}>Cancelar</Button>
            <Button type="submit">Registrar Horas</Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};


// --- Main Layout & Router ---

const Layout = ({ children }: { children?: React.ReactNode }) => {
  const location = useLocation();
  const navItems = [
    { label: 'Início', icon: LayoutDashboard, path: '/' },
    { label: 'Empresas', icon: Building2, path: '/companies' },
    { label: 'Contratos', icon: FileText, path: '/contracts' },
    { label: 'Usuários', icon: Users, path: '/users' },
    { label: 'Atendimentos', icon: Clock, path: '/appointments' },
  ];

  return (
    <div className="flex h-screen bg-gray-100 text-gray-900 font-sans overflow-hidden">
      {/* Sidebar */}
      <aside className="w-64 bg-slate-900 text-white flex-shrink-0 hidden md:flex flex-col border-r border-slate-800">
        <div className="h-16 flex items-center px-6 border-b border-slate-800">
          <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center mr-3">
            <span className="font-bold text-white">C5</span>
          </div>
          <span className="text-lg font-bold tracking-tight">Nexus</span>
        </div>
        
        <nav className="flex-1 px-4 py-6 space-y-1 overflow-y-auto">
          {navItems.map((item) => {
            const isActive = location.pathname === item.path;
            const Icon = item.icon;
            return (
              <Link
                key={item.path}
                to={item.path}
                className={`flex items-center px-3 py-2.5 text-sm font-medium rounded-md transition-colors ${
                  isActive 
                    ? 'bg-blue-600 text-white shadow-lg shadow-blue-900/50' 
                    : 'text-slate-400 hover:bg-slate-800 hover:text-white'
                }`}
              >
                <Icon className={`mr-3 w-5 h-5 ${isActive ? 'text-white' : 'text-slate-500'}`} />
                {item.label}
              </Link>
            );
          })}
        </nav>

        <div className="p-4 border-t border-slate-800">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-full bg-slate-700 flex items-center justify-center">
              <Users className="w-4 h-4 text-slate-300" />
            </div>
            <div>
              <p className="text-sm font-medium text-white">Administrador</p>
              <p className="text-xs text-slate-500">Gretti@c5.com</p>
            </div>
          </div>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 flex flex-col min-w-0 overflow-hidden">
        {/* Mobile Header (Placeholder for mobile toggle if needed) */}
        <div className="md:hidden bg-slate-900 text-white p-4 flex items-center justify-between">
          <span className="font-bold">Nexus</span>
        </div>

        <div className="flex-1 overflow-y-auto p-4 md:p-8">
          <div className="max-w-6xl mx-auto">
            {children}
          </div>
        </div>
      </main>
    </div>
  );
};

export default function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/companies" element={<CompaniesPage />} />
          <Route path="/contracts" element={<ContractsPage />} />
          <Route path="/users" element={<UsersPage />} />
          <Route path="/appointments" element={<AppointmentsPage />} />
        </Routes>
      </Layout>
    </Router>
  );
}