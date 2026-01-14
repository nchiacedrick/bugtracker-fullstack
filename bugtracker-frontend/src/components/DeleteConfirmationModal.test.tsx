import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import DeleteConfirmationModal from './DeleteConfirmationModal';

describe('DeleteConfirmationModal', () => {
    const mockOnClose = jest.fn();
    const mockOnConfirm = jest.fn();
    const bugTitle = 'Test Bug';

    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('should not render the modal when isOpen is false', () => {
        render(
            <DeleteConfirmationModal
                isOpen={false}
                onClose={mockOnClose}
                onConfirm={mockOnConfirm}
                bugTitle={bugTitle}
            />
        );
        expect(screen.queryByText('Confirm Delete')).not.toBeInTheDocument();
    });

    it('should render the modal when isOpen is true', () => {
        render(
            <DeleteConfirmationModal
                isOpen={true}
                onClose={mockOnClose}
                onConfirm={mockOnConfirm}
                bugTitle={bugTitle}
            />
        );
        expect(screen.getByText('Confirm Delete')).toBeInTheDocument();
        const dialog = screen.getByRole('dialog');
        expect(dialog).toHaveTextContent('Are you sure you want to delete bug');
        expect(dialog).toHaveTextContent(bugTitle);
        expect(dialog).toHaveTextContent('This action cannot be undone');
    });

    it('should call onClose when the Cancel button is clicked', () => {
        render(
            <DeleteConfirmationModal
                isOpen={true}
                onClose={mockOnClose}
                onConfirm={mockOnConfirm}
                bugTitle={bugTitle}
            />
        );
        fireEvent.click(screen.getByRole('button', { name: 'Cancel' }));
        expect(mockOnClose).toHaveBeenCalled();
    });

    it('should call onConfirm when the Delete button is clicked', () => {
        render(
            <DeleteConfirmationModal
                isOpen={true}
                onClose={mockOnClose}
                onConfirm={mockOnConfirm}
                bugTitle={bugTitle}
            />
        );
        fireEvent.click(screen.getByRole('button', { name: 'Delete' }));
        expect(mockOnConfirm).toHaveBeenCalled();
    });
}); 