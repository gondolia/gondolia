"use client";

import { useAuthStore } from "@/lib/stores/authStore";
import { Panel, PanelBody, PanelHeader } from "@/components/ui/Panel";

export default function DashboardPage() {
  const { user, company } = useAuthStore();

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
        Dashboard
      </h1>

      <div className="grid gap-6 md:grid-cols-2">
        <Panel>
          <PanelHeader>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              User Information
            </h2>
          </PanelHeader>
          <PanelBody>
            <dl className="space-y-3">
              <div>
                <dt className="text-sm text-gray-500 dark:text-gray-400">
                  Name
                </dt>
                <dd className="text-gray-900 dark:text-white font-medium">
                  {user?.displayName || `${user?.firstName} ${user?.lastName}`}
                </dd>
              </div>
              <div>
                <dt className="text-sm text-gray-500 dark:text-gray-400">
                  Email
                </dt>
                <dd className="text-gray-900 dark:text-white">
                  {user?.email}
                </dd>
              </div>
              <div>
                <dt className="text-sm text-gray-500 dark:text-gray-400">
                  Role
                </dt>
                <dd className="text-gray-900 dark:text-white">
                  {user?.role}
                </dd>
              </div>
            </dl>
          </PanelBody>
        </Panel>

        <Panel>
          <PanelHeader>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              Company Information
            </h2>
          </PanelHeader>
          <PanelBody>
            <dl className="space-y-3">
              <div>
                <dt className="text-sm text-gray-500 dark:text-gray-400">
                  Company Name
                </dt>
                <dd className="text-gray-900 dark:text-white font-medium">
                  {company?.name}
                </dd>
              </div>
              <div>
                <dt className="text-sm text-gray-500 dark:text-gray-400">
                  SAP Number
                </dt>
                <dd className="text-gray-900 dark:text-white">
                  {company?.sapNumber || "N/A"}
                </dd>
              </div>
              <div>
                <dt className="text-sm text-gray-500 dark:text-gray-400">
                  Status
                </dt>
                <dd>
                  <span
                    className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      company?.isActive
                        ? "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400"
                        : "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400"
                    }`}
                  >
                    {company?.isActive ? "Active" : "Inactive"}
                  </span>
                </dd>
              </div>
            </dl>
          </PanelBody>
        </Panel>
      </div>

      {user?.permissions && user.permissions.length > 0 && (
        <Panel>
          <PanelHeader>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              Permissions
            </h2>
          </PanelHeader>
          <PanelBody>
            <div className="flex flex-wrap gap-2">
              {user.permissions.map((permission) => (
                <span
                  key={permission}
                  className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary-100 text-primary-800 dark:bg-primary-900/30 dark:text-primary-400"
                >
                  {permission}
                </span>
              ))}
            </div>
          </PanelBody>
        </Panel>
      )}
    </div>
  );
}
