import { useEffect, useState, useCallback } from "react";
import { useRouter } from "next/router";
import { Bug } from "@/types/bug";
import Link from "next/link";
import CommentSection from "@/components/CommentSection";
import { Comment } from "@/types/comment";
import EditBugModal from "@/components/EditBugModal";
import { updateBug } from "@/api/bugs";
import DeleteConfirmationModal from "@/components/DeleteConfirmationModal";
import { deleteBug } from "@/api/bugs";
import { API_BASE_URL } from "@/config";

export default function BugDetail() {
  const router = useRouter();
  const { id } = router.query;
  const [bug, setBug] = useState<Bug | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);

  const fetchComments = useCallback(async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/bugs/${id}/comments`);
      if (!response.ok) {
        throw new Error("Failed to fetch comments");
      }
      const data = await response.json();
      setComments(data);
    } catch (error) {
      console.error("Error fetching comments:", error);
    }
  }, [id]);

  useEffect(() => {
    if (!id) return;

    const fetchData = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/api/bugs/${id}`);
        if (!response.ok) {
          throw new Error("Bug not found");
        }
        const data = await response.json();
        setBug(data);
        setLoading(false);
        await fetchComments();
      } catch (error) {
        console.error("Error fetching bug:", error);
        setError("Failed to fetch bug details");
        setLoading(false);
      }
    };

    fetchData();
  }, [id, fetchComments]);

  const handleEditBug = async (bugId: number, updatedBug: Partial<Bug>) => {
    try {
      await updateBug(bugId.toString(), updatedBug);
      const response = await fetch(`${API_BASE_URL}/api/bugs/${id}`);
      if (!response.ok) {
        throw new Error("Failed to fetch updated bug");
      }
      const data = await response.json();
      setBug(data);
      setIsEditModalOpen(false);
    } catch (error) {
      console.error("Failed to update bug:", error);
    }
  };

  const handleDelete = async () => {
    try {
      await deleteBug(bug!.id.toString());
      router.push({
        pathname: "/",
        query: {
          deletedBugTitle: bug!.title,
          showDeleteNotification: true,
        },
      });
    } catch (error) {
      console.error("Error deleting bug:", error);
    }
  };

  if (loading) return <div className="text-center p-4">Loading...</div>;
  if (error)
    return <div className="text-center p-4 text-red-500">Error: {error}</div>;
  if (!bug) return <div className="text-center p-4">Bug not found</div>;

  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-white shadow-lg">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <Link href="/" className="text-blue-500 hover:text-blue-700">
            ← Back to Bug List
          </Link>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="bg-white shadow-md rounded-lg p-6">
          <div className="flex justify-between items-center mb-4">
            <h1 className="text-2xl font-bold">{bug.title}</h1>
            <div className="flex gap-4">
              <button
                onClick={() => setIsEditModalOpen(true)}
                className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded"
              >
                Edit Bug
              </button>
              <button
                onClick={() => setIsDeleteModalOpen(true)}
                className="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded"
              >
                Delete Bug
              </button>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4 mb-6">
            <div>
              <p className="text-gray-600">ID</p>
              <p className="font-medium">{bug.id}</p>
            </div>
            <div>
              <p className="text-gray-600">Status</p>
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
            </div>
            <div>
              <p className="text-gray-600">Priority</p>
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
            </div>
          </div>

          <div className="mb-6">
            <h2 className="text-lg font-semibold mb-2">Description</h2>
            <p className="text-gray-700 whitespace-pre-wrap">
              {bug.description}
            </p>
          </div>

          <CommentSection
            bugId={bug.id}
            comments={comments || []}
            onCommentAdded={fetchComments}
          />

          {bug && (
            <EditBugModal
              isOpen={isEditModalOpen}
              onClose={() => setIsEditModalOpen(false)}
              onSubmit={handleEditBug}
              bug={bug}
            />
          )}

          {bug && (
            <DeleteConfirmationModal
              isOpen={isDeleteModalOpen}
              onClose={() => setIsDeleteModalOpen(false)}
              onConfirm={handleDelete}
              bugTitle={bug.title}
            />
          )}
        </div>
      </main>
    </div>
  );
}
