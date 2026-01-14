import React, { useEffect, useState } from "react";
import { getBugs, createBug, updateBug, deleteBug } from "../api/bugs";
import AddBugModal from "./AddBugModal";
import { Bug } from "../types/bug";
import Link from "next/link";
import EditBugModal from "./EditBugModal";
import { useRouter } from "next/router";
import DeleteConfirmationModal from "./DeleteConfirmationModal";
import Notification from "./Notification";
import { APP_VERSION } from "../config/app";
import Image from "next/image";

export default function BugList() {
  const router = useRouter();
  const { query, pathname, replace } = router;
  const [bugs, setBugs] = useState<Bug[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [selectedBug, setSelectedBug] = useState<Bug | null>(null);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [bugToDelete, setBugToDelete] = useState<Bug | null>(null);

  useEffect(() => {
    const fetchBugs = async () => {
      try {
        const data = await getBugs();
        setBugs(data);
        setError(null);
      } catch (error) {
        console.error("Error fetching bugs:", error);
        setError("Failed to fetch bugs");
      }
    };

    fetchBugs();
  }, []);

  const handleAddBug = async (newBug: Omit<Bug, "id" | "status">) => {
    try {
      const createdBug = await createBug({ ...newBug, status: "Open" });
      const updatedBugs = await getBugs();
      setBugs(updatedBugs);
      setIsModalOpen(false);

      replace({
        pathname,
        query: {
          createdBugTitle: createdBug.title,
          showCreateNotification: true,
        },
      });
    } catch (error) {
      console.error("Failed to create bug:", error);
    }
  };

  const handleEditClick = (bug: Bug) => {
    setSelectedBug(bug);
    setIsEditModalOpen(true);
  };

  const handleEditBug = async (bugId: number, updatedBug: Partial<Bug>) => {
    try {
      await updateBug(bugId.toString(), updatedBug);
      const updatedBugs = await getBugs();
      setBugs(updatedBugs);
      setIsEditModalOpen(false);
      setSelectedBug(null);
    } catch (error) {
      console.error("Failed to update bug:", error);
    }
  };

  const handleDeleteClick = (bug: Bug) => {
    setBugToDelete(bug);
    setIsDeleteModalOpen(true);
  };

  const handleConfirmDelete = async () => {
    if (!bugToDelete) return;

    try {
      await deleteBug(bugToDelete.id.toString());
      const updatedBugs = await getBugs();
      setBugs(updatedBugs);
      setIsDeleteModalOpen(false);
      setBugToDelete(null);

      replace({
        pathname,
        query: {
          deletedBugTitle: bugToDelete.title,
          showDeleteNotification: true,
        },
      });
    } catch (error) {
      console.error("Failed to delete bug:", error);
    }
  };

  if (error)
    return <div className="text-center p-4 text-red-500">Error: {error}</div>;

  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-white shadow-lg">
        <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
          <div className="flex items-center">
            <Image
              src="/bugTracker_Logo.png"
              alt="Bug Tracker Logo"
              width={40}
              height={40}
              className="mr-3"
              priority
            />
            <h1 className="text-2xl font-bold text-gray-800">
              Bug Tracker Pro
            </h1>
          </div>
          <span className="text-gray-600">v{APP_VERSION}</span>
        </div>
      </nav>

      {query.showDeleteNotification && (
        <Notification
          message={`Successfully deleted bug "${query.deletedBugTitle}"`}
          onClose={() => {
            replace({
              pathname,
              query: {},
            });
          }}
          type="success"
        />
      )}

      {query.showCreateNotification && (
        <Notification
          message={`Successfully created bug "${query.createdBugTitle}"`}
          onClose={() => {
            replace({
              pathname,
              query: {},
            });
          }}
          type="success"
        />
      )}

      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-xl font-semibold text-gray-700">All Bugs</h2>
          <button
            onClick={() => setIsModalOpen(true)}
            className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-md"
          >
            Add New Bug
          </button>
        </div>

        <div className="bg-white shadow-md rounded-lg">
          <table className="min-w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  ID
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Title
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Priority
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {(bugs || []).map((bug) => (
                <tr key={bug.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    <Link
                      href={`/bugs/${bug.id}`}
                      className="text-blue-600 hover:text-blue-900"
                    >
                      {bug.id}
                    </Link>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    <Link
                      href={`/bugs/${bug.id}`}
                      className="text-blue-600 hover:text-blue-900"
                    >
                      {bug.title}
                    </Link>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span
                      className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full
                                            ${
                                              bug.status === "Open"
                                                ? "bg-red-100 text-red-800"
                                                : bug.status === "In Progress"
                                                ? "bg-yellow-100 text-yellow-800"
                                                : "bg-green-100 text-green-800"
                                            }`}
                    >
                      {bug.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span
                      className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full
                                            ${
                                              bug.priority === "High"
                                                ? "bg-red-100 text-red-800"
                                                : bug.priority === "Medium"
                                                ? "bg-yellow-100 text-yellow-800"
                                                : "bg-green-100 text-green-800"
                                            }`}
                    >
                      {bug.priority}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button
                      onClick={() => handleEditClick(bug)}
                      className="text-blue-600 hover:text-blue-900 mr-4"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDeleteClick(bug)}
                      className="text-red-600 hover:text-red-900"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </main>

      <AddBugModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleAddBug}
      />

      {selectedBug && (
        <EditBugModal
          isOpen={isEditModalOpen}
          onClose={() => {
            setIsEditModalOpen(false);
            setSelectedBug(null);
          }}
          onSubmit={handleEditBug}
          bug={selectedBug}
        />
      )}

      {bugToDelete && (
        <DeleteConfirmationModal
          isOpen={isDeleteModalOpen}
          onClose={() => {
            setIsDeleteModalOpen(false);
            setBugToDelete(null);
          }}
          onConfirm={handleConfirmDelete}
          bugTitle={bugToDelete.title}
        />
      )}
    </div>
  );
}
