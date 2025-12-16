'use client'

import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api, Contract, Company, ContractReport } from '@/lib/api'
import { format } from 'date-fns'
import Modal from '@/components/Modal'

export default function ContractsPage() {
  const queryClient = useQueryClient()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingContract, setEditingContract] = useState<Contract | null>(null)
  const [selectedReport, setSelectedReport] = useState<ContractReport | null>(null)

  const { data: contracts = [], isLoading } = useQuery({
    queryKey: ['contracts'],
    queryFn: api.contracts.list,
  })

  const { data: companies = [] } = useQuery({
    queryKey: ['companies'],
    queryFn: api.companies.list,
  })

  const createMutation = useMutation({
    mutationFn: api.contracts.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['contracts'] })
      setIsModalOpen(false)
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Contract> }) =>
      api.contracts.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['contracts'] })
      setIsModalOpen(false)
      setEditingContract(null)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: api.contracts.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['contracts'] })
    },
  })

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const formData = new FormData(e.currentTarget)
    const data = {
      companyId: formData.get('companyId') as string,
      title: formData.get('title') as string,
      totalHours: parseFloat(formData.get('totalHours') as string),
      startDate: formData.get('startDate') as string,
      endDate: formData.get('endDate') as string,
      type: formData.get('type') as 'monthly' | 'project' | 'retainer',
    }

    if (editingContract) {
      updateMutation.mutate({ id: editingContract.id, data })
    } else {
      createMutation.mutate(data)
    }
  }

  const openEditModal = (contract: Contract) => {
    setEditingContract(contract)
    setIsModalOpen(true)
  }

  const closeModal = () => {
    setIsModalOpen(false)
    setEditingContract(null)
  }

  const viewReport = async (contractId: string) => {
    const report = await api.contracts.report(contractId)
    setSelectedReport(report)
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Contracts</h1>
        <button
          onClick={() => setIsModalOpen(true)}
          className="btn btn-primary"
        >
          Add Contract
        </button>
      </div>

      {isLoading ? (
        <p className="text-gray-500">Loading...</p>
      ) : contracts.length === 0 ? (
        <div className="card text-center py-12">
          <p className="text-gray-500 mb-4">No contracts yet</p>
          <button
            onClick={() => setIsModalOpen(true)}
            className="btn btn-primary"
          >
            Add your first contract
          </button>
        </div>
      ) : (
        <div className="card overflow-hidden p-0">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="table-header">Title</th>
                <th className="table-header">Company</th>
                <th className="table-header">Hours</th>
                <th className="table-header">Period</th>
                <th className="table-header">Status</th>
                <th className="table-header text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {contracts.map((contract: Contract) => (
                <tr key={contract.id} className="hover:bg-gray-50">
                  <td className="table-cell font-medium">{contract.title}</td>
                  <td className="table-cell">{contract.companyName}</td>
                  <td className="table-cell">{contract.totalHours}h</td>
                  <td className="table-cell">
                    {format(new Date(contract.startDate), 'MMM d')} -{' '}
                    {format(new Date(contract.endDate), 'MMM d, yyyy')}
                  </td>
                  <td className="table-cell">
                    <span
                      className={`px-2 py-1 text-xs rounded-full ${
                        contract.isActive
                          ? 'bg-green-100 text-green-700'
                          : 'bg-gray-100 text-gray-600'
                      }`}
                    >
                      {contract.isActive ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td className="table-cell text-right">
                    <button
                      onClick={() => viewReport(contract.id)}
                      className="text-green-600 hover:text-green-800 mr-4"
                    >
                      Report
                    </button>
                    <button
                      onClick={() => openEditModal(contract)}
                      className="text-primary-600 hover:text-primary-800 mr-4"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => {
                        if (confirm('Are you sure you want to delete this contract?')) {
                          deleteMutation.mutate(contract.id)
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
        title={editingContract ? 'Edit Contract' : 'Add Contract'}
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="label">Contract Title</label>
            <input
              name="title"
              type="text"
              required
              defaultValue={editingContract?.title}
              className="input"
              placeholder="e.g., Monthly Support Contract"
            />
          </div>
          <div>
            <label className="label">Company</label>
            <select
              name="companyId"
              required
              defaultValue={editingContract?.companyId}
              className="input"
            >
              <option value="">Select a company</option>
              {companies.map((company: Company) => (
                <option key={company.id} value={company.id}>
                  {company.name}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="label">Total Hours</label>
            <input
              name="totalHours"
              type="number"
              step="0.5"
              required
              defaultValue={editingContract?.totalHours}
              className="input"
              placeholder="e.g., 160"
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="label">Start Date</label>
              <input
                name="startDate"
                type="date"
                required
                defaultValue={
                  editingContract
                    ? format(new Date(editingContract.startDate), 'yyyy-MM-dd')
                    : ''
                }
                className="input"
              />
            </div>
            <div>
              <label className="label">End Date</label>
              <input
                name="endDate"
                type="date"
                required
                defaultValue={
                  editingContract
                    ? format(new Date(editingContract.endDate), 'yyyy-MM-dd')
                    : ''
                }
                className="input"
              />
            </div>
          </div>
          <div>
            <label className="label">Contract Type</label>
            <select
              name="type"
              required
              defaultValue={editingContract?.type || 'monthly'}
              className="input"
            >
              <option value="monthly">Monthly</option>
              <option value="project">Project</option>
              <option value="retainer">Retainer</option>
            </select>
          </div>
          <div className="flex justify-end gap-3 pt-4">
            <button type="button" onClick={closeModal} className="btn btn-secondary">
              Cancel
            </button>
            <button type="submit" className="btn btn-primary">
              {editingContract ? 'Save Changes' : 'Create Contract'}
            </button>
          </div>
        </form>
      </Modal>

      <Modal
        isOpen={!!selectedReport}
        onClose={() => setSelectedReport(null)}
        title="Contract Report"
      >
        {selectedReport && (
          <div className="space-y-4">
            <div className="border-b pb-4">
              <h3 className="font-semibold text-lg">{selectedReport.contractTitle}</h3>
              <p className="text-gray-500">{selectedReport.companyName}</p>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-gray-50 p-4 rounded-lg">
                <p className="text-sm text-gray-500">Total Hours</p>
                <p className="text-2xl font-bold">{selectedReport.totalHours}h</p>
              </div>
              <div className="bg-blue-50 p-4 rounded-lg">
                <p className="text-sm text-blue-600">Consumed Hours</p>
                <p className="text-2xl font-bold text-blue-700">
                  {selectedReport.consumedHours}h
                </p>
              </div>
              <div
                className={`p-4 rounded-lg ${
                  selectedReport.isOverBudget ? 'bg-red-50' : 'bg-green-50'
                }`}
              >
                <p
                  className={`text-sm ${
                    selectedReport.isOverBudget ? 'text-red-600' : 'text-green-600'
                  }`}
                >
                  Remaining Balance
                </p>
                <p
                  className={`text-2xl font-bold ${
                    selectedReport.isOverBudget ? 'text-red-700' : 'text-green-700'
                  }`}
                >
                  {selectedReport.remainingBalance}h
                </p>
              </div>
              <div className="bg-purple-50 p-4 rounded-lg">
                <p className="text-sm text-purple-600">Days Remaining</p>
                <p className="text-2xl font-bold text-purple-700">
                  {selectedReport.daysRemaining}
                </p>
              </div>
            </div>
            <div className="bg-yellow-50 p-4 rounded-lg">
              <p className="text-sm text-yellow-700 mb-2">Ideal Monthly Burn</p>
              <p className="text-xl font-bold text-yellow-800">
                {selectedReport.idealMonthlyBurn}h / month
              </p>
              <p className="text-sm text-yellow-600 mt-1">
                Current burn rate: {selectedReport.actualBurnRate}h / month
              </p>
            </div>
            <div className="flex justify-end pt-4">
              <button
                onClick={() => setSelectedReport(null)}
                className="btn btn-secondary"
              >
                Close
              </button>
            </div>
          </div>
        )}
      </Modal>
    </div>
  )
}
