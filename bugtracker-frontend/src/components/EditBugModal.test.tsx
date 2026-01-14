import { render, screen, fireEvent } from '@testing-library/react';
import EditBugModal from './EditBugModal';

describe('EditBugModal', () => {
    const mockBug = {
        id: 1,
        title: 'Test Bug',
        description: 'Test Description',
        priority: 'Medium' as const,
        status: 'Open' as const,
        createdAt: new Date().toISOString()
    };

    const mockProps = {
        isOpen: true,
        onClose: jest.fn(),
        onSubmit: jest.fn(),
        bug: mockBug
    };

    beforeEach(() => {
        jest.clearAllMocks();
    });

    test('renders nothing when isOpen is false', () => {
        render(<EditBugModal {...mockProps} isOpen={false} />);
        expect(screen.queryByText('Edit Bug')).not.toBeInTheDocument();
    });

    test('renders modal content when isOpen is true', () => {
        render(<EditBugModal {...mockProps} />);
        expect(screen.getByText('Edit Bug')).toBeInTheDocument();
        expect(screen.getByDisplayValue('Test Bug')).toBeInTheDocument();
        expect(screen.getByDisplayValue('Test Description')).toBeInTheDocument();
    });

    test('calls onClose when Cancel button is clicked', () => {
        render(<EditBugModal {...mockProps} />);
        fireEvent.click(screen.getByText('Cancel'));
        expect(mockProps.onClose).toHaveBeenCalledTimes(1);
    });

    test('updates form fields when input values change', () => {
        render(<EditBugModal {...mockProps} />);
        
        const titleInput = screen.getByDisplayValue('Test Bug');
        fireEvent.change(titleInput, { target: { value: 'Updated Bug Title' } });
        expect(screen.getByDisplayValue('Updated Bug Title')).toBeInTheDocument();

        const descriptionInput = screen.getByDisplayValue('Test Description');
        fireEvent.change(descriptionInput, { target: { value: 'Updated Description' } });
        expect(screen.getByDisplayValue('Updated Description')).toBeInTheDocument();
    });

    test('calls onSubmit with updated values when form is submitted', () => {
        render(<EditBugModal {...mockProps} />);
        
        const titleInput = screen.getByDisplayValue('Test Bug');
        fireEvent.change(titleInput, { target: { value: 'Updated Bug Title' } });

        const form = screen.getByTestId('edit-bug-form');
        fireEvent.submit(form);

        expect(mockProps.onSubmit).toHaveBeenCalledWith(mockBug.id, expect.objectContaining({
            title: 'Updated Bug Title',
            description: 'Test Description',
            priority: 'Medium',
            status: 'Open'
        }));
        expect(mockProps.onClose).toHaveBeenCalled();
    });

    test('updates priority and status when select values change', () => {
        render(<EditBugModal {...mockProps} />);
        
        const prioritySelect = screen.getByDisplayValue('Medium');
        fireEvent.change(prioritySelect, { target: { value: 'High' } });
        expect(screen.getByDisplayValue('High')).toBeInTheDocument();

        const statusSelect = screen.getByDisplayValue('Open');
        fireEvent.change(statusSelect, { target: { value: 'In Progress' } });
        expect(screen.getByDisplayValue('In Progress')).toBeInTheDocument();
    });
});
