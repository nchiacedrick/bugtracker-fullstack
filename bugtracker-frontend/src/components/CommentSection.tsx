import { useState } from "react";
import { Comment } from "@/types/comment";
import { API_BASE_URL } from "@/config";

interface CommentSectionProps {
  bugId: number;
  comments: Comment[];
  onCommentAdded: () => void;
}

export default function CommentSection({
  bugId,
  comments = [],
  onCommentAdded,
}: CommentSectionProps) {
  const [author, setAuthor] = useState("");
  const [content, setContent] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      const response = await fetch(
        `${API_BASE_URL}/api/bugs/${bugId}/comments`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ author, content }),
        }
      );

      if (!response.ok) {
        throw new Error("Failed to add comment");
      }

      setAuthor("");
      setContent("");
      onCommentAdded();
    } catch (error) {
      console.error("Error adding comment:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="mt-8">
      <h2 className="text-xl font-semibold mb-4">Comments</h2>

      <div className="mb-6">
        <form
          data-testid="comment-form"
          onSubmit={handleSubmit}
          className="space-y-4"
        >
          <div>
            <label
              htmlFor="author"
              className="block text-sm font-medium text-gray-700"
            >
              Your Name
            </label>
            <input
              data-testid="comment-author"
              type="text"
              id="author"
              value={author}
              onChange={(e) => setAuthor(e.target.value)}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              required
            />
          </div>
          <div>
            <label
              htmlFor="content"
              className="block text-sm font-medium text-gray-700"
            >
              Comment
            </label>
            <textarea
              data-testid="comment-content"
              id="content"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              rows={3}
              required
            />
          </div>
          <button
            type="submit"
            disabled={isSubmitting || !author || !content}
            className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 disabled:opacity-50"
          >
            {isSubmitting ? "Adding..." : "Add Comment"}
          </button>
        </form>
      </div>

      <div className="space-y-4">
        {comments.length === 0 ? (
          <p>No comments yet.</p>
        ) : (
          comments.map((comment: Comment) => (
            <div key={comment.id} className="bg-gray-50 p-4 rounded-lg">
              <div className="flex justify-between items-start">
                <span className="font-medium">{comment.author}</span>
                <span className="text-sm text-gray-500">
                  {new Date(comment.createdAt).toLocaleString("en-US", {
                    dateStyle: "short",
                    timeStyle: "medium",
                  })}
                </span>
              </div>
              <p className="mt-2 text-gray-700">{comment.content}</p>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
