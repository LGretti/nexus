import { Company, Contract, User, Appointment, ApiError } from '../types';

const API_BASE_URL = 'http://localhost:8080/api';

async function fetchClient<T>(endpoint: string, options?: RequestInit): Promise<T> {
  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    // Handle 204 No Content (often used for successful deletes or empty lists)
    if (response.status === 204) {
      return [] as unknown as T;
    }

    const text = await response.text();
    let data;
    try {
       // Handle empty body
       data = text ? JSON.parse(text) : null;
    } catch (e) {
       console.error("Failed to parse JSON", e);
       throw new Error("Invalid JSON response from server");
    }

    if (!response.ok) {
      const errorMsg = (data as ApiError)?.error || `Error ${response.status}: ${response.statusText}`;
      throw new Error(errorMsg);
    }

    // Critical Fix: Go backend returns 'null' for nil slices. 
    // We convert this to [] to prevent frontend crashes on .map()
    if (data === null) {
      return [] as unknown as T;
    }

    return data as T;
  } catch (error: any) {
    // Handle network errors (Failed to fetch) specifically
    if (error.message === 'Failed to fetch') {
      throw new Error('Cannot connect to API. Is the Go backend running at port 8080?');
    }
    throw error;
  }
}

// Companies
export const getCompanies = () => fetchClient<Company[]>('/companies');
export const createCompany = (company: Omit<Company, 'id'>) => 
  fetchClient<Company>('/companies', { method: 'POST', body: JSON.stringify(company) });
export const updateCompany = (id: number, company: Omit<Company, 'id'>) => 
  fetchClient<Company>(`/companies/${id}`, { method: 'PUT', body: JSON.stringify(company) });
export const deleteCompany = (id: number) => 
  fetchClient<void>(`/companies/${id}`, { method: 'DELETE' });

// Contracts
export const getContracts = () => fetchClient<Contract[]>('/contracts');
// Fallback: If /contracts/{id} doesn't exist on backend, we might rely on the list, 
// but per spec we add a method to fetch a single contract if the API supports it.
export const getContractById = (id: number) => fetchClient<Contract>(`/contracts/${id}`); 
export const createContract = (contract: Omit<Contract, 'id' | 'companyName'>) => 
  fetchClient<Contract>('/contracts', { method: 'POST', body: JSON.stringify(contract) });
export const getContractAppointments = (id: number) => 
  fetchClient<Appointment[]>(`/contracts/${id}/appointments`);

// Users
export const getUsers = () => fetchClient<User[]>('/users');
export const createUser = (user: Omit<User, 'id'>) => 
  fetchClient<User>('/users', { method: 'POST', body: JSON.stringify(user) });

// Appointments
export const getAppointments = () => fetchClient<Appointment[]>('/appointments');
export const createAppointment = (appointment: Omit<Appointment, 'id' | 'contractTitle' | 'userName' | 'totalHours'>) => 
  fetchClient<Appointment>('/appointments', { method: 'POST', body: JSON.stringify(appointment) });
export const deleteAppointment = (id: number) => 
  fetchClient<void>(`/appointments/${id}`, { method: 'DELETE' });
