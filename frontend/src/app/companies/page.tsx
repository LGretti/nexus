'use client'

import { useEffect, useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api, Company } from '@/lib/api'
import Modal from '@/components/Modal'

export default function CompaniesPage() {
  const queryClient = useQueryClient()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingCompany, setEditingCompany] = useState<Company | null>(null)

  const { data: companies = [], isLoading } = useQuery({
    queryKey: ['companies'],
    queryFn: api.companies.list,
  })

  const createMutation = useMutation({
    mutationFn: api.companies.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['companies'] })
      setIsModalOpen(false)
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Company> }) =>
      api.companies.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['companies'] })
      setIsModalOpen(false)
      setEditingCompany(null)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: api.companies.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['companies'] })
    },
  })

  useEffect(() => {
  const errorMsg = createMutation.error?.message || updateMutation.error?.message
  
  if (errorMsg) {
    // Verifica se a mensagem fala de "CNPJ" (convertendo para minúsculo pra garantir)
    if (errorMsg.toLowerCase().includes('cnpj')) {
      // Busca o input pelo 'name="cnpj"' e dá o foco
      const cnpjInput = document.querySelector('input[name="cnpj"]') as HTMLInputElement
      if (cnpjInput) {
        cnpjInput.focus()
        cnpjInput.select() // Opcional: Seleciona o texto para o usuário já digitar por cima
      }
    }
    
    // Bônus: Se o erro for de Email duplicado (futuramente), já deixa pronto:
    if (errorMsg.toLowerCase().includes('email')) {
      const emailInput = document.querySelector('input[name="email"]') as HTMLInputElement
      emailInput?.focus()
    }
  }
}, [createMutation.error, updateMutation.error])

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const formData = new FormData(e.currentTarget)
    const data = {
      name: formData.get('name') as string,
      email: formData.get('email') as string,
      cnpj: formData.get('cnpj') as string,
    }

    if (editingCompany) {
      updateMutation.mutate({ id: editingCompany.id, data })
    } else {
      createMutation.mutate(data)
    }
  }

  const openEditModal = (company: Company) => {
    setEditingCompany(company)
    setIsModalOpen(true)
  }

  const closeModal = () => {
    setIsModalOpen(false)
    setEditingCompany(null)
    createMutation.reset()
    updateMutation.reset()
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Companies</h1>
        <button
          onClick={() => setIsModalOpen(true)}
          className="btn btn-primary"
        >
          Add Company
        </button>
      </div>

      {isLoading ? (
        <p className="text-gray-500">Loading...</p>
      ) : companies.length === 0 ? (
        <div className="card text-center py-12">
          <p className="text-gray-500 mb-4">No companies yet</p>
          <button
            onClick={() => setIsModalOpen(true)}
            className="btn btn-primary"
          >
            Add your first company
          </button>
        </div>
      ) : (
        <div className="card overflow-hidden p-0">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="table-header">Name</th>
                <th className="table-header">Email</th>
                <th className="table-header">CNPJ</th>
                <th className="table-header text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {companies.map((company: Company) => (
                <tr key={company.id} className="hover:bg-gray-50">
                  <td className="table-cell font-medium">{company.name}</td>
                  <td className="table-cell">{company.email}</td>
                  <td className="table-cell">{company.cnpj}</td>
                  <td className="table-cell text-right">
                    <button
                      onClick={() => openEditModal(company)}
                      className="text-primary-600 hover:text-primary-800 mr-4"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => {
                        if (confirm('Are you sure you want to delete this company?')) {
                          deleteMutation.mutate(company.id)
                        }
                      }}
                      className="text-red-600 hover:text-red-800"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <Modal
        isOpen={isModalOpen}
        onClose={closeModal}
        title={editingCompany ? 'Edit Company' : 'Add Company'}
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          {(createMutation.isError || updateMutation.isError) && (
            <div className="bg-red-50 border-l-4 border-red-500 p-4 mb-4">
              <div className="flex">
                <div className="ml-3">
                  <p className="text-sm text-red-700">
                    {/* Pega a mensagem de erro do Create OU do Update */}
                    {createMutation.error?.message || updateMutation.error?.message || "Ocorreu um erro ao salvar."}
                  </p>
                </div>
              </div>
            </div>
          )}
          {/* =================================== */}
          <div>
            <label className="label">Company Name</label>
            <input
              name="name"
              type="text"
              required
              defaultValue={editingCompany?.name}
              className="input"
              placeholder="e.g., Ademicon"
            />
          </div>
          <div>
            <label className="label">Email</label>
            <input
              name="email"
              type="email"
              required
              defaultValue={editingCompany?.email}
              className="input"
              placeholder="e.g., contact@company.com"
            />
          </div>
          <div>
            <label className="label">CNPJ</label>
            <input
              name="cnpj"
              type="text"
              required
              defaultValue={editingCompany?.cnpj}
              className="input"
              placeholder="00.000.000/0000-00"
            />
          </div>
          <div className="flex justify-end gap-3 pt-4">
            <button type="button" onClick={closeModal} className="btn btn-secondary">
              Cancel
            </button>
            <button type="submit" className="btn btn-primary">
              {editingCompany ? 'Save Changes' : 'Create Company'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
