import { useEffect, useState } from "react";

interface NotificationProps {
  message: string;
  onClose: () => void;
  type?: "success" | "error";
}

export default function Notification({
  message,
  onClose,
  type = "success",
}: NotificationProps) {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    setIsVisible(true);
    const timer = setTimeout(() => {
      setIsVisible(false);
    }, 4500);
    const closeTimer = setTimeout(onClose, 5000);

    return () => {
      clearTimeout(timer);
      clearTimeout(closeTimer);
    };
  }, [onClose]);

  const bgColor =
    type === "success"
      ? "bg-green-100 border-green-400 text-green-700"
      : "bg-red-100 border-red-400 text-red-700";
  const iconColor = type === "success" ? "text-green-500" : "text-red-500";

  return (
    <div
      className={`max-w-7xl mx-auto px-4 py-4 transition-opacity duration-500 ${
        isVisible ? "opacity-100" : "opacity-0"
      }`}
    >
      <div
        className={`${bgColor} border px-4 py-3 rounded relative`}
        role="alert"
      >
        <span className="block sm:inline">{message}</span>
        <span className="absolute top-0 bottom-0 right-0 px-4 py-3">
          <button onClick={onClose}>
            <svg
              className={`fill-current h-6 w-6 ${iconColor}`}
              role="button"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 20 20"
            >
              <title>Close</title>
              <path d="M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10 11.819l-2.651 3.029a1.2 1.2 0 1 1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1 1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697 1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z" />
            </svg>
          </button>
        </span>
      </div>
    </div>
  );
}
