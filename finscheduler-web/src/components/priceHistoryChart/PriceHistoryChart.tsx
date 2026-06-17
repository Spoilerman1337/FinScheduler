import {Box, Card} from '@chakra-ui/react';
import {ArrowDown, ArrowUp, Minus} from 'lucide-react';
import {useId, useMemo, useState} from 'react';
import {
    PRICE_HISTORY_ACTUAL_STROKE,
    PRICE_HISTORY_FORECAST_STROKE,
    PRICE_HISTORY_PRICE_DECREASE_COLOR,
    PRICE_HISTORY_PRICE_INCREASE_COLOR,
    PRICE_HISTORY_PRICE_UNCHANGED_COLOR,
    PRICE_HISTORY_TOOLTIP_WIDTH_PX,
} from './shared.ts';
import type {
    PriceHistoryChartMouseState,
    PriceHistoryChartPoint,
    PriceHistoryLegendItem,
    PriceHistoryMetricItem,
    PriceHistoryPoint,
    PriceHistoryTooltipContent,
    PriceHistoryTooltipEntry,
    PriceHistoryTooltipPosition,
} from './shared.ts';
import PriceHistoryChartCanvas from './subcomponents/PriceHistoryChartCanvas.tsx';
import PriceHistoryEmptyState from './subcomponents/PriceHistoryEmptyState.tsx';
import PriceHistoryHeader from './subcomponents/PriceHistoryHeader.tsx';
import PriceHistoryMetrics from './subcomponents/PriceHistoryMetrics.tsx';

interface PriceHistoryChartProps {
    points?: PriceHistoryPoint[] | null;
    forecastPoints?: PriceHistoryPoint[] | null;
}

interface ChartSummary {
    current: number;
    min: number;
    max: number;
}

interface TooltipValue {
    color?: string;
    value: number;
    absoluteChange: number | null;
    percentChange: number | null;
    changeColor: string;
}

interface ActiveTooltipState {
    coordinate: {
        x: number;
        y: number;
    };
    index: number;
}

const tooltipDateFormatter = new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
});

const axisDateFormatter = new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: 'short',
});

const moneyFormatter = new Intl.NumberFormat('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
});

const percentFormatter = new Intl.NumberFormat('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
});

function toSortedSeries(points?: PriceHistoryPoint[] | null): PriceHistoryPoint[] {
    return [...(points ?? [])].sort((left, right) => {
        return new Date(left.point).getTime() - new Date(right.point).getTime();
    });
}

function buildForecastSeries(
    actualPoints: PriceHistoryPoint[],
    forecastPoints: PriceHistoryPoint[],
) {
    if (actualPoints.length === 0 || forecastPoints.length === 0) {
        return forecastPoints;
    }

    const lastActualPoint = actualPoints.at(-1)!;
    const firstForecastPoint = forecastPoints[0];

    if (lastActualPoint.point === firstForecastPoint.point) {
        return forecastPoints;
    }

    return [lastActualPoint, ...forecastPoints];
}

function buildChartData(
    actualPoints: PriceHistoryPoint[],
    forecastPoints: PriceHistoryPoint[],
): PriceHistoryChartPoint[] {
    const actualMap = new Map(actualPoints.map((point) => [point.point, point]));
    const forecastMap = new Map(forecastPoints.map((point) => [point.point, point]));
    const points = Array.from(new Set([...actualMap.keys(), ...forecastMap.keys()])).sort(
        (left, right) => {
            return new Date(left).getTime() - new Date(right).getTime();
        },
    );

    return points.map((point) => {
        const actualPoint = actualMap.get(point);
        const forecastPoint = forecastMap.get(point);

        return {
            point,
            actualValue: actualPoint?.value ?? null,
            actualAbsoluteChange: actualPoint?.absoluteChange ?? null,
            actualPercentChange: actualPoint?.percentChange ?? null,
            forecastValue: forecastPoint?.value ?? null,
            forecastAbsoluteChange: forecastPoint?.absoluteChange ?? null,
            forecastPercentChange: forecastPoint?.percentChange ?? null,
        };
    });
}

function buildSummary(points: PriceHistoryPoint[]): ChartSummary | null {
    if (points.length === 0) {
        return null;
    }

    const values = points.map(({value}) => value);

    return {
        current: points.at(-1)!.value,
        min: Math.min(...values),
        max: Math.max(...values),
    };
}

function formatMoney(value: number) {
    return `${moneyFormatter.format(value)} ₽`;
}

function formatAxisMoney(value: number) {
    const normalizedValue = Number.isInteger(value) ? value.toFixed(0) : value.toFixed(2);

    return `${normalizedValue.replace('.', ',')} ₽`;
}

function formatAxisDate(value: string) {
    return axisDateFormatter.format(new Date(value));
}

function formatTooltipDate(value: string) {
    return tooltipDateFormatter.format(new Date(value));
}

function buildTooltipValues(point: PriceHistoryChartPoint): TooltipValue[] {
    const values: TooltipValue[] = [];

    if (point.actualValue !== null) {
        values.push({
            color: PRICE_HISTORY_ACTUAL_STROKE,
            value: point.actualValue,
            absoluteChange: point.actualAbsoluteChange,
            percentChange: point.actualPercentChange,
            changeColor: getPriceChangeColor(point.actualAbsoluteChange),
        });
    }

    if (point.forecastValue !== null) {
        const forecastValue = {
            color: PRICE_HISTORY_FORECAST_STROKE,
            value: point.forecastValue,
            absoluteChange: point.forecastAbsoluteChange,
            percentChange: point.forecastPercentChange,
            changeColor: getPriceChangeColor(point.forecastAbsoluteChange),
        };
        const alreadyIncluded = values.some((entry) => {
            return entry.color === forecastValue.color && entry.value === forecastValue.value;
        });

        if (!alreadyIncluded) {
            values.push(forecastValue);
        }
    }

    return values;
}

function buildTooltipEntries(values: TooltipValue[]): PriceHistoryTooltipEntry[] {
    return values.map((entry) => ({
        key: `${entry.color}-${entry.value}`,
        value: formatMoney(entry.value),
        valueColor: entry.color ?? 'fg',
        changeSummary:
            entry.absoluteChange !== null
                ? formatPriceChangeSummary(entry.absoluteChange, entry.percentChange)
                : undefined,
        changeColor: entry.absoluteChange !== null ? entry.changeColor : undefined,
        changeIcon:
            entry.absoluteChange !== null ? getPriceChangeIcon(entry.absoluteChange) : undefined,
    }));
}

function buildTooltipPosition(coordinate: {x: number; y: number}): PriceHistoryTooltipPosition {
    const halfTooltipWidth = PRICE_HISTORY_TOOLTIP_WIDTH_PX / 2;

    return {
        left: `clamp(8px, calc(${coordinate.x}px - ${halfTooltipWidth}px), calc(100% - ${PRICE_HISTORY_TOOLTIP_WIDTH_PX}px - 8px))`,
        top: `${coordinate.y}px`,
        transform: coordinate.y < 72 ? 'translateY(14px)' : 'translateY(calc(-100% - 14px))',
    };
}

function getPriceChangeColor(change: number | null) {
    if (change === null) {
        return 'var(--chakra-colors-fg-subtle)';
    }

    if (change > 0) {
        return PRICE_HISTORY_PRICE_INCREASE_COLOR;
    }

    if (change < 0) {
        return PRICE_HISTORY_PRICE_DECREASE_COLOR;
    }

    return PRICE_HISTORY_PRICE_UNCHANGED_COLOR;
}

function formatAbsoluteMoney(value: number) {
    return formatMoney(Math.abs(value));
}

function formatSignedPercent(value: number) {
    if (value === 0) {
        return `${percentFormatter.format(0)}%`;
    }

    const sign = value > 0 ? '+' : '-';

    return `${sign}${percentFormatter.format(Math.abs(value))}%`;
}

function getPriceChangeIcon(change: number | null) {
    if (change === null || change === 0) {
        return Minus;
    }

    return change > 0 ? ArrowUp : ArrowDown;
}

function formatPriceChangeSummary(absoluteChange: number, percentChange: number | null) {
    const absoluteMoneyPart = `(${formatAbsoluteMoney(absoluteChange)})`;

    if (percentChange === null) {
        return absoluteMoneyPart;
    }

    return `${formatSignedPercent(percentChange)} ${absoluteMoneyPart}`;
}

export default function PriceHistoryChart({points, forecastPoints}: PriceHistoryChartProps) {
    const [activeTooltip, setActiveTooltip] = useState<ActiveTooltipState | null>(null);
    const actualPoints = useMemo(() => toSortedSeries(points), [points]);
    const preparedForecastPoints = useMemo(() => {
        return buildForecastSeries(actualPoints, toSortedSeries(forecastPoints));
    }, [actualPoints, forecastPoints]);
    const chartData = useMemo(() => {
        return buildChartData(actualPoints, preparedForecastPoints);
    }, [actualPoints, preparedForecastPoints]);
    const summary = useMemo(() => buildSummary(actualPoints), [actualPoints]);
    const hasForecast = preparedForecastPoints.length > 0;
    const legends: PriceHistoryLegendItem[] = [
        {label: 'Факт', color: PRICE_HISTORY_ACTUAL_STROKE},
        ...(hasForecast ? [{label: 'Прогноз', color: PRICE_HISTORY_FORECAST_STROKE}] : []),
    ];
    const metricItems: PriceHistoryMetricItem[] = summary
        ? [
              {
                  label: 'Минимум',
                  value: formatMoney(summary.min),
                  accentColor: 'neon.blue',
              },
              {
                  label: 'Текущая',
                  value: formatMoney(summary.current),
                  accentColor: 'neon.purple',
              },
              {
                  label: 'Максимум',
                  value: formatMoney(summary.max),
                  accentColor: 'neon.yellow',
              },
          ]
        : [];
    const activeTooltipPoint =
        activeTooltip && chartData[activeTooltip.index] ? chartData[activeTooltip.index] : null;
    const activeTooltipValues = useMemo(() => {
        if (!activeTooltipPoint) {
            return [];
        }

        return buildTooltipValues(activeTooltipPoint);
    }, [activeTooltipPoint]);
    const activeTooltipPosition = activeTooltip
        ? buildTooltipPosition(activeTooltip.coordinate)
        : null;
    const activeTooltipContent = useMemo<PriceHistoryTooltipContent | null>(() => {
        if (!activeTooltipPosition || !activeTooltipPoint || activeTooltipValues.length === 0) {
            return null;
        }

        return {
            label: formatTooltipDate(activeTooltipPoint.point),
            entries: buildTooltipEntries(activeTooltipValues),
            position: activeTooltipPosition,
        };
    }, [activeTooltipPoint, activeTooltipPosition, activeTooltipValues]);

    const handleChartMouseMove = (state: PriceHistoryChartMouseState) => {
        const index = Number(state.activeTooltipIndex);
        const coordinateX =
            typeof state.activeCoordinate?.x === 'number' ? state.activeCoordinate.x : state.chartX;
        const coordinateY =
            typeof state.activeCoordinate?.y === 'number' ? state.activeCoordinate.y : state.chartY;

        if (
            !Number.isFinite(index) ||
            typeof coordinateX !== 'number' ||
            typeof coordinateY !== 'number'
        ) {
            setActiveTooltip(null);
            return;
        }

        setActiveTooltip({
            coordinate: {
                x: coordinateX,
                y: coordinateY,
            },
            index,
        });
    };

    const gradientId = useId();
    const forecastGradientId = useId();

    return (
        <Card.Root overflow="hidden">
            <Box
                position="absolute"
                inset="0"
                pointerEvents="none"
                bg="radial-gradient(circle at top right, rgba(143, 120, 255, 0.18), transparent 45%)"
            />
            <PriceHistoryHeader
                title="История цены"
                description="Фактические точки по дням."
                legends={legends}
            />
            <Card.Body position="relative" gap={4}>
                {summary ? (
                    <>
                        <PriceHistoryMetrics items={metricItems} />
                        <PriceHistoryChartCanvas
                            chartLabel={
                                hasForecast
                                    ? 'График истории и прогноза цены'
                                    : 'График истории цены'
                            }
                            data={chartData}
                            hasForecast={hasForecast}
                            activeTooltip={activeTooltipContent}
                            onMouseMove={handleChartMouseMove}
                            onMouseLeave={() => setActiveTooltip(null)}
                            gradientId={gradientId}
                            forecastGradientId={forecastGradientId}
                            formatAxisDate={formatAxisDate}
                            formatAxisMoney={formatAxisMoney}
                        />
                    </>
                ) : (
                    <PriceHistoryEmptyState
                        title="История цены появится после первой фиксации изменения стоимости."
                        description="Пока здесь будут отображаться только фактические точки. Когда появятся прогнозы, их можно будет дорисовать второй серией без смены карточки."
                    />
                )}
            </Card.Body>
        </Card.Root>
    );
}
