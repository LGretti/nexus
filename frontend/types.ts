export interface Company {
  id: number;
  name: string;
  cnpj: string;
  email: string;
}

export interface User {
  id: number;
  name: string;
  email: string;
  role: 'admin' | 'consultant';
}

export interface Contract {
  id: number;
  companyId: number;
  companyName?: string; // Read-only from API
  title: string;
  contractType: string;
  totalHours: number;
  startDate: string;
  endDate: string;
  isActive: boolean;
}

export interface Appointment {
  id: number;
  contractId: number;
  userId: number;
  contractTitle?: string; // Read-only from API
  userName?: string; // Read-only from API
  startTime: string;
  endTime: string;
  description: string;
  totalHours?: number; // Read-only from API
  durationSeconds?: number;
}

export interface ApiError {
  error: string;
}
