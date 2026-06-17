import {fireEvent, screen, waitFor} from '@testing-library/react';
import {cloneElement, isValidElement, type ReactElement} from 'react';
import {describe, expect, it, vi} from 'vitest';
import {renderWithProviders} from '../../test/render.tsx';
import PriceHistoryChart from './PriceHistoryChart.tsx';

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

describe('PriceHistoryChart', () => {
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
        const coordinateX = Number(hoveredDot?.getAttribute('cx') ?? 0);
        const coordinateY = Number(hoveredDot?.getAttribute('cy') ?? 0);

        expect(hoveredDot).toBeTruthy();

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
        fireEvent.mouseLeave(hoveredDot);

        // Assert
        await waitFor(() => {
            expect(screen.queryByTestId('price-history-tooltip')).not.toBeInTheDocument();
        });
    });
});
