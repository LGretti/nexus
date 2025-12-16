'use client'

import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api, Appointment, Contract, User } from '@/lib/api'
import { format } from 'date-fns'
import Modal from '@/components/Modal'

export default function AppointmentsPage() {
  const queryClient = useQueryClient()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingAppointment, setEditingAppointment] = useState<Appointment | null>(null)

  const { data: appointments = [], isLoading } = useQuery({
    queryKey: ['appointments'],
    queryFn: api.appointments.list,
  })

  const { data: contracts = [] } = useQuery({
    queryKey: ['contracts'],
    queryFn: api.contracts.list,
  })

  const { data: users = [] } = useQuery({
    queryKey: ['users'],
    queryFn: api.users.list,
  })

  const activeContracts = contracts.filter((c: Contract) => c.isActive)

  const createMutation = useMutation({
    mutationFn: api.appointments.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['appointments'] })
      setIsModalOpen(false)
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Appointment> }) =>
      api.appointments.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['appointments'] })
      setIsModalOpen(false)
      setEditingAppointment(null)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: api.appointments.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['appointments'] })
    },
  })

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const formData = new FormData(e.currentTarget)
    const data = {
      contractId: formData.get('contractId') as string,
      userId: formData.get('userId') as string,
      date: formData.get('date') as string,
      hours: parseFloat(formData.get('hours') as string),
      description: formData.get('description') as string,
    }

    if (editingAppointment) {
      updateMutation.mutate({ id: editingAppointment.id, data })
    } else {
      createMutation.mutate(data)
    }
  }

  const openEditModal = (appointment: Appointment) => {
    setEditingAppointment(appointment)
    setIsModalOpen(true)
  }

  const closeModal = () => {
    setIsModalOpen(false)
    setEditingAppointment(null)
  }

  const totalHoursLogged = appointments.reduce(
    (sum: number, a: Appointment) => sum + a.hours,
    0
  )

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Time Tracking</h1>
          <p className="text-gray-500 mt-1">
            Total hours logged: <span className="font-semibold">{totalHoursLogged}h</span>
          </p>
        </div>
        <button
          onClick={() => setIsModalOpen(true)}
          className="btn btn-primary"
        >
          Log Time
        </button>
      </div>

      {isLoading ? (
        <p className="text-gray-500">Loading...</p>
      ) : appointments.length === 0 ? (
        <div className="card text-center py-12">
          <p className="text-gray-500 mb-4">No time entries yet</p>
          <button
            onClick={() => setIsModalOpen(true)}
            className="btn btn-primary"
          >
            Log your first time entry
          </button>
        </div>
      ) : (
        <div className="card overflow-hidden p-0">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="table-header">Date</th>
                <th className="table-header">Contract</th>
                <th className="table-header">Consultant</th>
                <th className="table-header">Hours</th>
                <th className="table-header">Description</th>
                <th className="table-header text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {appointments.map((appointment: Appointment) => (
                <tr key={appointment.id} className="hover:bg-gray-50">
                  <td className="table-cell">
                    {format(new Date(appointment.date), 'MMM d, yyyy')}
                  </td>
                  <td className="table-cell font-medium">
                    {appointment.contractName}
                  </td>
                  <td className="table-cell">{appointment.userName}</td>
                  <td className="table-cell">
                    <span className="font-semibold text-primary-600">
                      {appointment.hours}h
                    </span>
                  </td>
                  <td className="table-cell max-w-xs truncate">
                    {appointment.description}
                  </td>
                  <td className="table-cell text-right">
                    <button
                      onClick={() => openEditModal(appointment)}
                      className="text-primary-600 hover:text-primary-800 mr-4"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => {
                        if (confirm('Are you sure you want to delete this entry?')) {
                          deleteMutation.mutate(appointment.id)
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
        title={editingAppointment ? 'Edit Time Entry' : 'Log Time'}
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="label">Contract</label>
            <select
              name="contractId"
              required
              defaultValue={editingAppointment?.contractId}
              className="input"
            >
              <option value="">Select an active contract</option>
              {activeContracts.map((contract: Contract) => (
                <option key={contract.id} value={contract.id}>
                  {contract.title} ({contract.companyName})
                </option>
              ))}
            </select>
            {activeContracts.length === 0 && (
              <p className="text-sm text-red-500 mt-1">
                No active contracts available
              </p>
            )}
          </div>
          <div>
            <label className="label">Consultant</label>
            <select
              name="userId"
              required
              defaultValue={editingAppointment?.userId}
              className="input"
            >
              <option value="">Select a user</option>
              {users.map((user: User) => (
                <option key={user.id} value={user.id}>
                  {user.name} ({user.role})
                </option>
              ))}
            </select>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="label">Date</label>
              <input
                name="date"
                type="date"
                required
                defaultValue={
                  editingAppointment
                    ? format(new Date(editingAppointment.date), 'yyyy-MM-dd')
                    : format(new Date(), 'yyyy-MM-dd')
                }
                className="input"
              />
            </div>
            <div>
              <label className="label">Hours</label>
              <input
                name="hours"
                type="number"
                step="0.25"
                min="0.25"
                max="24"
                required
                defaultValue={editingAppointment?.hours || 1}
                className="input"
                placeholder="e.g., 4"
              />
            </div>
          </div>
          <div>
            <label className="label">Description</label>
            <textarea
              name="description"
              rows={3}
              required
              defaultValue={editingAppointment?.description}
              className="input"
              placeholder="What did you work on?"
            />
          </div>
          <div className="flex justify-end gap-3 pt-4">
            <button type="button" onClick={closeModal} className="btn btn-secondary">
              Cancel
            </button>
            <button
              type="submit"
              className="btn btn-primary"
              disabled={activeContracts.length === 0}
            >
              {editingAppointment ? 'Save Changes' : 'Log Time'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
