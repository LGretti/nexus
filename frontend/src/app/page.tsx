'use client'

import { useQuery } from '@tanstack/react-query'
import { api, Contract } from '@/lib/api'
import { format } from 'date-fns'

export default function Dashboard() {
  const { data: contracts = [] } = useQuery({
    queryKey: ['contracts'],
    queryFn: api.contracts.list,
  })

  const { data: companies = [] } = useQuery({
    queryKey: ['companies'],
    queryFn: api.companies.list,
  })

  const { data: users = [] } = useQuery({
    queryKey: ['users'],
    queryFn: api.users.list,
  })

  const { data: appointments = [] } = useQuery({
    queryKey: ['appointments'],
    queryFn: api.appointments.list,
  })

  const activeContracts = contracts.filter((c: Contract) => c.isActive)

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatCard
          title="Companies"
          value={companies.length}
          icon="🏢"
          color="blue"
        />
        <StatCard
          title="Active Contracts"
          value={activeContracts.length}
          icon="📄"
          color="green"
        />
        <StatCard
          title="Users"
          value={users?.length}
          icon="👥"
          color="purple"
        />
        <StatCard
          title="Time Entries"
          value={appointments?.length}
          icon="⏱️"
          color="orange"
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            Active Contracts
          </h2>
          {activeContracts.length === 0 ? (
            <p className="text-gray-500">No active contracts</p>
          ) : (
            <div className="space-y-4">
              {activeContracts.slice(0, 5).map((contract: Contract) => (
                <ContractCard key={contract.id} contract={contract} />
              ))}
            </div>
          )}
        </div>

        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            Recent Time Entries
          </h2>
          {appointments?.length === 0 ? (
            <p className="text-gray-500">No time entries yet</p>
          ) : (
            <div className="space-y-3">
              {appointments?.slice(0, 5).map((appt: any) => (
                <div
                  key={appt.id}
                  className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
                >
                  <div>
                    <p className="font-medium text-gray-900">
                      {appt.contractName || 'Contract'}
                    </p>
                    <p className="text-sm text-gray-500">
                      {appt.userName || 'User'} -{' '}
                      {format(new Date(appt.date), 'MMM d, yyyy')}
                    </p>
                  </div>
                  <span className="font-semibold text-primary-600">
                    {appt.hours}h
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

function StatCard({
  title,
  value,
  icon,
  color,
}: {
  title: string
  value: number
  icon: string
  color: 'blue' | 'green' | 'purple' | 'orange'
}) {
  const colorClasses = {
    blue: 'bg-blue-50 text-blue-600',
    green: 'bg-green-50 text-green-600',
    purple: 'bg-purple-50 text-purple-600',
    orange: 'bg-orange-50 text-orange-600',
  }

  return (
    <div className="card">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-gray-500">{title}</p>
          <p className="text-3xl font-bold text-gray-900 mt-1">{value}</p>
        </div>
        <div
          className={`w-12 h-12 rounded-full flex items-center justify-center text-2xl ${colorClasses[color]}`}
        >
          {icon}
        </div>
      </div>
    </div>
  )
}

function ContractCard({ contract }: { contract: Contract }) {
  return (
    <div className="p-4 bg-gray-50 rounded-lg">
      <div className="flex items-center justify-between mb-2">
        <h3 className="font-medium text-gray-900">{contract.title}</h3>
        <span
          className={`px-2 py-1 text-xs rounded-full ${
            contract.isActive
              ? 'bg-green-100 text-green-700'
              : 'bg-gray-100 text-gray-600'
          }`}
        >
          {contract.isActive ? 'Active' : 'Inactive'}
        </span>
      </div>
      <p className="text-sm text-gray-500 mb-2">
        {contract.companyName || 'Company'}
      </p>
      <div className="flex items-center justify-between text-sm">
        <span className="text-gray-500">
          {format(new Date(contract.startDate), 'MMM d')} -{' '}
          {format(new Date(contract.endDate), 'MMM d, yyyy')}
        </span>
        <span className="font-semibold text-primary-600">
          {contract.totalHours}h
        </span>
      </div>
    </div>
  )
}
