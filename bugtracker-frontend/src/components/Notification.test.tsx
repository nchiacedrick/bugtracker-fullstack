import { render, screen, act, fireEvent } from "@testing-library/react";
import Notification from "./Notification";

describe("Notification", () => {
  beforeEach(() => {
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it("should render with the provided message", () => {
    render(<Notification message="Test notification" onClose={jest.fn()} />);
    expect(screen.getByText("Test notification")).toBeInTheDocument();
  });

  it("should call onClose after 5 seconds", () => {
    const onClose = jest.fn();
    render(<Notification message="Test notification" onClose={onClose} />);

    expect(onClose).not.toHaveBeenCalled();

    act(() => {
      jest.advanceTimersByTime(5000);
    });

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it("should call onClose when close button is clicked", () => {
    const onClose = jest.fn();
    render(<Notification message="Test notification" onClose={onClose} />);

    const closeButton = screen.getByRole("button", { name: /close/i });
    fireEvent.click(closeButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it("should apply success styles by default", () => {
    render(<Notification message="Success notification" onClose={jest.fn()} />);

    const notification = screen.getByRole("alert", { hidden: true });
    expect(notification).toHaveClass(
      "bg-green-100",
      "border-green-400",
      "text-green-700"
    );
  });

  it("should apply error styles when type is error", () => {
    render(
      <Notification
        message="Error notification"
        onClose={jest.fn()}
        type="error"
      />
    );

    const notification = screen.getByRole("alert", { hidden: true });
    expect(notification).toHaveClass(
      "bg-red-100",
      "border-red-400",
      "text-red-700"
    );
  });

  it("should fade out before closing", () => {
    render(<Notification message="Test notification" onClose={jest.fn()} />);

    const notification = screen.getByRole("alert", {
      hidden: true,
    }).parentElement;
    expect(notification).toHaveClass("opacity-100");

    act(() => {
      jest.advanceTimersByTime(4500);
    });

    expect(notification).toHaveClass("opacity-0");
  });

  it("should cleanup timers on unmount", () => {
    const onClose = jest.fn();
    const { unmount } = render(
      <Notification message="Test notification" onClose={onClose} />
    );

    unmount();

    act(() => {
      jest.advanceTimersByTime(5000);
    });

    expect(onClose).not.toHaveBeenCalled();
  });
});
