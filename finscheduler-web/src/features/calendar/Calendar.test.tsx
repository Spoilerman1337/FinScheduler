import {act, fireEvent, screen, waitFor, within} from '@testing-library/react';
import {afterEach, beforeEach, describe, expect, it, vi} from 'vitest';
import {renderWithProviders} from '../../test/render.tsx';
import Calendar from './Calendar.tsx';

describe('Calendar page', () => {
    beforeEach(() => {
        vi.useFakeTimers();
        vi.setSystemTime(new Date(2026, 6, 15, 12, 0, 0));
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    it('renders a full month view and shows earliest-time chips without no-time badges', () => {
        // Arrange
        renderWithProviders(<Calendar />);
        const detailsRail = screen.getByLabelText('Сведения по выбранному дню');
        const currentDayButton = screen.getByRole('button', {
            name: 'Открыть 15 июля 2026, 2 события',
        });
        const secondMarkedDayButton = screen.getByRole('button', {
            name: 'Открыть 18 июля 2026, 2 события',
        });
        const multipleTimesChip = within(currentDayButton).getByLabelText(
            'У события Проверка категорий самое раннее время 10:00, есть еще слоты',
        );
        const eventTitle = multipleTimesChip.parentElement
            ?.previousElementSibling as HTMLElement | null;
        const noTimeIndicator = within(currentDayButton).queryByRole('img', {
            name: 'У события Лимиты на неделю время не указано',
        });

        // Assert
        expect(screen.getByText('Июль 2026')).toBeInTheDocument();
        expect(screen.getByLabelText('Календарная сетка')).toBeInTheDocument();
        expect(within(detailsRail).getByText('Среда, 15 июля')).toBeInTheDocument();
        expect(within(detailsRail).getByText('Проверка категорий')).toBeInTheDocument();
        expect(within(detailsRail).getByText('Лимиты на неделю')).toBeInTheDocument();
        expect(within(detailsRail).getByText('10:00')).toBeInTheDocument();
        expect(within(detailsRail).getByText('13:30')).toBeInTheDocument();
        expect(within(detailsRail).getByText('18:00')).toBeInTheDocument();
        expect(within(detailsRail).queryByText('Время не указано')).not.toBeInTheDocument();
        expect(currentDayButton).toHaveStyle({
            height: '7rem',
            minHeight: '7rem',
            maxHeight: '7rem',
        });
        expect(secondMarkedDayButton).toHaveStyle({
            height: '7rem',
            minHeight: '7rem',
            maxHeight: '7rem',
        });
        expect(eventTitle).toHaveStyle({
            whiteSpace: 'nowrap',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
        });
        expect(within(multipleTimesChip).getByText('10:00+')).toBeInTheDocument();
        expect(within(currentDayButton).queryByText(/^2 /)).not.toBeInTheDocument();
        expect(within(secondMarkedDayButton).getByText('13:00')).toBeInTheDocument();
        expect(within(secondMarkedDayButton).queryByText('15:30')).not.toBeInTheDocument();
        expect(multipleTimesChip).toHaveStyle({
            height: '1.25rem',
        });
        expect(noTimeIndicator).not.toBeInTheDocument();
    });

    it('updates the details rail when another marked day is selected', async () => {
        // Arrange
        renderWithProviders(<Calendar />);
        const detailsRail = screen.getByLabelText('Сведения по выбранному дню');
        vi.useRealTimers();

        // Act
        await act(async () => {
            fireEvent.click(screen.getByRole('button', {name: 'Открыть 18 июля 2026, 2 события'}));
        });

        // Assert
        await waitFor(() => {
            expect(within(detailsRail).getByText('Суббота, 18 июля')).toBeInTheDocument();
        });
        expect(within(detailsRail).getByText('Запланировать перевод')).toBeInTheDocument();
        expect(within(detailsRail).getByText('Подтвердить зачисление')).toBeInTheDocument();
        expect(within(detailsRail).getByText('13:00')).toBeInTheDocument();
        expect(within(detailsRail).getByText('15:30')).toBeInTheDocument();
        expect(
            screen.getByRole('button', {name: 'Открыть 18 июля 2026, 2 события'}),
        ).toHaveAttribute('aria-pressed', 'true');
    });

    it('shows the empty-state rail for an unmarked day', async () => {
        // Arrange
        renderWithProviders(<Calendar />);
        const detailsRail = screen.getByLabelText('Сведения по выбранному дню');
        vi.useRealTimers();

        // Act
        await act(async () => {
            fireEvent.click(
                screen.getByRole('button', {name: 'Открыть 16 июля 2026, без событий'}),
            );
        });

        // Assert
        await waitFor(() => {
            expect(within(detailsRail).getByText('Четверг, 16 июля')).toBeInTheDocument();
        });
        expect(within(detailsRail).getByText('Событий на эту дату нет')).toBeInTheDocument();
        expect(
            within(detailsRail).getByText('Эта дата пока не привязана ни к одному событию.'),
        ).toBeInTheDocument();
    });
});
