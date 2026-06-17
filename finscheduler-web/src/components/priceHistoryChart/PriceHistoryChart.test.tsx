import {fireEvent, screen, waitFor, within} from '@testing-library/react';
import {cloneElement, isValidElement, type ReactElement} from 'react';
import {describe, expect, it, vi} from 'vitest';
import {renderWithProviders} from '../../test/render.tsx';
import PriceHistoryChart, {resolveActiveTooltipIndex} from './PriceHistoryChart.tsx';

interface ResponsiveContainerChildProps {
    width?: number;
    height?: number;
}

vi.mock('recharts', async () => {
    const actual = await vi.importActual<typeof import('recharts')>('recharts');

    return {
        ...actual,
        ResponsiveContainer: ({
            children,
        }: {
            children: ReactElement<ResponsiveContainerChildProps>;
        }) => {
            if (!isValidElement<ResponsiveContainerChildProps>(children)) {
                return null;
            }

            return cloneElement(children, {
                width: 360,
                height: 220,
            });
        },
    };
});

function renderChart() {
    return renderWithProviders(
        <PriceHistoryChart
            points={[
                {
                    point: '2026-01-12T00:00:00Z',
                    value: 139.9,
                    absoluteChange: null,
                    percentChange: null,
                },
                {
                    point: '2026-02-18T00:00:00Z',
                    value: 149.5,
                    absoluteChange: 9.6,
                    percentChange: 6.861,
                },
                {
                    point: '2026-03-20T00:00:00Z',
                    value: 160,
                    absoluteChange: 10.5,
                    percentChange: 7.023,
                },
            ]}
        />,
    );
}

function renderChartWithForecast() {
    return renderWithProviders(
        <PriceHistoryChart
            points={[
                {
                    point: '2026-01-12T00:00:00Z',
                    value: 139.9,
                    absoluteChange: null,
                    percentChange: null,
                },
                {
                    point: '2026-02-18T00:00:00Z',
                    value: 149.5,
                    absoluteChange: 9.6,
                    percentChange: 6.861,
                },
                {
                    point: '2026-03-20T00:00:00Z',
                    value: 160,
                    absoluteChange: 10.5,
                    percentChange: 7.023,
                },
            ]}
            forecastPoints={[
                {
                    point: '2026-04-20T00:00:00Z',
                    value: 165.75,
                },
                {
                    point: '2026-05-20T00:00:00Z',
                    value: 171.5,
                },
            ]}
        />,
    );
}

describe('PriceHistoryChart', () => {
    it('returns null when recharts does not provide an active tooltip index', () => {
        // Arrange
        const pointsCount = 5;

        // Act
        const missingIndex = resolveActiveTooltipIndex(null, pointsCount);
        const emptyIndex = resolveActiveTooltipIndex('', pointsCount);

        // Assert
        expect(missingIndex).toBeNull();
        expect(emptyIndex).toBeNull();
    });

    it('shows the hovered point price and date in the tooltip', async () => {
        // Arrange
        const {container} = renderChart();
        const dots = container.querySelectorAll('.recharts-line-dots circle');
        const hoveredDot = dots[1];

        expect(hoveredDot).toBeTruthy();

        // Act
        fireEvent.mouseEnter(hoveredDot);
        fireEvent.mouseMove(hoveredDot, {
            clientX: Number(hoveredDot?.getAttribute('cx') ?? 0),
            clientY: Number(hoveredDot?.getAttribute('cy') ?? 0),
        });

        // Assert
        await waitFor(() => {
            expect(screen.getByText(/149,50/)).toBeInTheDocument();
        });
        expect(screen.getByText(/\+6,86% \(9,60/)).toBeInTheDocument();
        expect(
            screen.getByText((content) => {
                return content.includes('18') && content.includes('2026');
            }),
        ).toBeInTheDocument();
    });

    it('hides the tooltip after the pointer leaves the chart point', async () => {
        // Arrange
        const {container} = renderChart();
        const dots = container.querySelectorAll('.recharts-line-dots circle');
        const hoveredDot = dots[dots.length - 1];
        const chartWrapper = container.querySelector('.recharts-wrapper');
        const coordinateX = Number(hoveredDot?.getAttribute('cx') ?? 0);
        const coordinateY = Number(hoveredDot?.getAttribute('cy') ?? 0);

        expect(hoveredDot).toBeTruthy();
        expect(chartWrapper).toBeTruthy();

        // Act
        fireEvent.mouseEnter(hoveredDot);
        fireEvent.mouseMove(hoveredDot, {
            clientX: coordinateX,
            clientY: coordinateY,
        });

        // Assert
        await waitFor(() => {
            expect(screen.getByTestId('price-history-tooltip')).toBeInTheDocument();
        });
        fireEvent.mouseLeave(chartWrapper!);

        // Assert
        await waitFor(() => {
            expect(screen.queryByTestId('price-history-tooltip')).not.toBeInTheDocument();
        });
    });

    it('renders a stretchable time range brush for the chart', () => {
        // Arrange
        const {container} = renderChartWithForecast();

        // Act
        const brush = container.querySelector('.recharts-brush');

        // Assert
        expect(brush).toBeTruthy();
    });

    it('renders the forecast series and shows the forecast value in the tooltip', async () => {
        // Arrange
        const {container} = renderChartWithForecast();
        const dots = container.querySelectorAll('.recharts-line-dots circle');
        const hoveredDot = dots[dots.length - 1];

        expect(dots.length).toBeGreaterThan(3);
        expect(hoveredDot).toBeTruthy();

        // Act
        fireEvent.mouseEnter(hoveredDot);
        fireEvent.mouseMove(hoveredDot, {
            clientX: Number(hoveredDot?.getAttribute('cx') ?? 0),
            clientY: Number(hoveredDot?.getAttribute('cy') ?? 0),
        });

        // Assert
        await waitFor(() => {
            expect(screen.getByText(/171,50/)).toBeInTheDocument();
        });
        expect(
            within(screen.getByTestId('price-history-tooltip')).getByText('Прогноз'),
        ).toBeInTheDocument();
    });

    it('does not show forecast details on the last actual point used to connect the series', async () => {
        // Arrange
        const {container} = renderChartWithForecast();
        const dots = container.querySelectorAll('.recharts-line-dots circle');
        const actualLastDot = dots[2];

        expect(actualLastDot).toBeTruthy();

        // Act
        fireEvent.mouseEnter(actualLastDot);
        fireEvent.mouseMove(actualLastDot, {
            clientX: Number(actualLastDot?.getAttribute('cx') ?? 0),
            clientY: Number(actualLastDot?.getAttribute('cy') ?? 0),
        });

        // Assert
        await waitFor(() => {
            expect(screen.getByTestId('price-history-tooltip')).toBeInTheDocument();
        });
        expect(
            within(screen.getByTestId('price-history-tooltip')).getAllByText(/160,00/),
        ).toHaveLength(1);
    });
});
