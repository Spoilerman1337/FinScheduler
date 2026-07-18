import {fireEvent, screen} from '@testing-library/react';
import {afterEach, beforeEach, describe, expect, it, vi} from 'vitest';
import {renderWithProviders} from '../../test/render.tsx';
import Events from './Events.tsx';

describe('Events page', () => {
    beforeEach(() => {
        vi.useFakeTimers();
        vi.setSystemTime(new Date(2026, 6, 18, 12, 0, 0));
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    it('renders only the name and date range filters together with the events table', () => {
        // Arrange
        renderWithProviders(<Events />);

        // Assert
        expect(screen.getByRole('columnheader', {name: 'Событие'})).toBeInTheDocument();
        expect(screen.getByRole('columnheader', {name: 'Дата'})).toBeInTheDocument();
        expect(screen.getByRole('columnheader', {name: 'Время'})).toBeInTheDocument();
        expect(screen.getByPlaceholderText('Поиск по наименованию...')).toBeInTheDocument();
        expect(screen.getByRole('button', {name: 'Дата'})).toBeInTheDocument();
        expect(screen.getAllByText('Запланировать перевод').length).toBeGreaterThan(0);
        expect(screen.getAllByText('Проверка категорий').length).toBeGreaterThan(0);
        expect(screen.queryByRole('columnheader', {name: 'Триггеры'})).not.toBeInTheDocument();
        expect(screen.queryByText('Триггеры:')).not.toBeInTheDocument();
        expect(screen.queryByRole('button', {name: 'Сегодня'})).not.toBeInTheDocument();
    });

    it('filters events only by the event name', () => {
        // Arrange
        renderWithProviders(<Events />);

        // Act
        fireEvent.change(screen.getByPlaceholderText('Поиск по наименованию...'), {
            target: {value: 'перевод'},
        });

        // Assert
        expect(screen.getAllByText('Запланировать перевод').length).toBeGreaterThan(0);
        expect(screen.queryByText('Подтвердить зачисление')).not.toBeInTheDocument();
        expect(screen.queryByText('Проверка категорий')).not.toBeInTheDocument();
    });
});